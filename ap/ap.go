// Package ap implements the Spotify accesspoint protocol: the Diffie-Hellman +
// Shannon handshake over TCP and the encrypted packet channel. It is used to
// request per-file AES audio keys, which are not available over the REST API.
package ap

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"log/slog"
	"net"
	"slices"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	"google.golang.org/protobuf/proto"
	"pyrorhythm.dev/libspot/audio/dh"
	pb "pyrorhythm.dev/libspot/gen/spotify"
)

const pongAckInterval = 120 * time.Second

// GetAddressFunc resolves an accesspoint address ("host:port"). It is called on
// every (re)connection so a fresh accesspoint can be selected.
type GetAddressFunc func(ctx context.Context) string

type AccesspointLoginError struct {
	Message *pb.APLoginFailed
}

func (e *AccesspointLoginError) Error() string {
	return fmt.Sprintf("accesspoint login failed: %s %v",
		e.Message.GetErrorCode().String(), e.Message.GetErrorDescription())
}

type Accesspoint struct {
	log *slog.Logger

	addr GetAddressFunc

	nonce    []byte
	deviceId string

	dh *dh.DiffieHellman

	conn    net.Conn
	encConn *shannonConn

	stop              bool
	pongAckTickerStop chan struct{}
	recvLoopStop      chan struct{}
	recvLoopOnce      sync.Once
	recvChans         map[PacketType][]chan Packet
	recvChansLock     sync.RWMutex
	lastPongAck       time.Time
	lastPongAckLock   sync.Mutex

	// connMu is held for writing when reconnecting and for reading when sending
	// packets or accessing welcome.
	connMu  sync.RWMutex
	welcome *pb.APWelcome
}

func NewAccesspoint(log *slog.Logger, addr GetAddressFunc, deviceId string) *Accesspoint {
	return &Accesspoint{log: log, addr: addr, deviceId: deviceId, recvChans: make(map[PacketType][]chan Packet)}
}

func (ap *Accesspoint) init(ctx context.Context) (err error) {
	// read 16 nonce bytes
	ap.nonce = make([]byte, 16)
	if _, err = rand.Read(ap.nonce); err != nil {
		return fmt.Errorf("failed reading random nonce: %w", err)
	}

	// init diffiehellman parameters
	if ap.dh, err = dh.NewDiffieHellman(); err != nil {
		return fmt.Errorf("failed initializing diffiehellman: %w", err)
	}

	// close previous connection if any
	if ap.conn != nil {
		_ = ap.conn.Close()
		ap.conn = nil
	}

	// open connection to accesspoint
	var dialer net.Dialer
	for attempts := 1; ; attempts++ {
		ctx_, cancel := context.WithTimeout(ctx, time.Second*30)
		addr := ap.addr(ctx_)
		conn, err := dialer.DialContext(ctx_, "tcp", addr)
		cancel()
		if err == nil {
			ap.conn = conn
			ap.log.Debug("connected to accesspoint", "addr", addr)
			return nil
		} else if attempts >= 6 {
			return fmt.Errorf("failed to connect to AP %v: %w", addr, err)
		}
		ap.log.Warn("failed to connect to AP, retrying with a different AP", "addr", addr, "error", err)
	}
}

func (ap *Accesspoint) ConnectSpotifyToken(ctx context.Context, username, token string) error {
	return ap.Connect(ctx, pb.LoginCredentials_builder{
		Typ:      pb.AuthenticationType_AUTHENTICATION_SPOTIFY_TOKEN.Enum(),
		Username: proto.String(username),
		AuthData: []byte(token),
	}.Build())
}

func (ap *Accesspoint) Connect(ctx context.Context, creds *pb.LoginCredentials) error {
	ap.connMu.Lock()
	defer ap.connMu.Unlock()

	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		if err := ap.connect(ctx, creds); err != nil {
			ap.log.Warn("failed connecting to accesspoint, retrying", "error", err)
			return struct{}{}, err
		}
		return struct{}{}, nil
	}, backoff.WithBackOff(backoff.NewConstantBackOff(500*time.Millisecond)), backoff.WithMaxTries(5))
	return err
}

func (ap *Accesspoint) connect(ctx context.Context, creds *pb.LoginCredentials) error {
	ap.recvLoopStop = make(chan struct{}, 1)
	ap.pongAckTickerStop = make(chan struct{}, 1)

	if err := ap.init(ctx); err != nil {
		return err
	}

	if deadline, ok := ctx.Deadline(); ok {
		_ = ap.conn.SetDeadline(deadline)
		defer func() { _ = ap.conn.SetDeadline(time.Time{}) }()
	}

	// perform key exchange with diffiehellman
	exchangeData, err := ap.performKeyExchange()
	if err != nil {
		return fmt.Errorf("failed performing keyexchange: %w", err)
	}

	// solve challenge and complete connection
	if err := ap.solveChallenge(exchangeData); err != nil {
		return fmt.Errorf("failed solving challenge: %w", err)
	}

	// do authentication with credentials
	if err := ap.authenticate(ctx, creds); err != nil {
		return fmt.Errorf("failed authenticating: %w", err)
	}

	return nil
}

func (ap *Accesspoint) Close() {
	ap.connMu.Lock()
	defer ap.connMu.Unlock()

	ap.stop = true

	if ap.conn == nil {
		return
	}

	ap.recvLoopStop <- struct{}{}
	ap.pongAckTickerStop <- struct{}{}
	_ = ap.conn.Close()
}

func (ap *Accesspoint) Send(ctx context.Context, pktType PacketType, payload []byte) error {
	ap.connMu.RLock()
	defer ap.connMu.RUnlock()
	return ap.encConn.sendPacket(ctx, pktType, payload)
}

func (ap *Accesspoint) Receive(types ...PacketType) <-chan Packet {
	ch := make(chan Packet)
	ap.recvChansLock.Lock()
	for _, type_ := range types {
		ap.recvChans[type_] = append(ap.recvChans[type_], ch)
	}
	ap.recvChansLock.Unlock()

	ap.startReceiving()

	return ch
}

func (ap *Accesspoint) startReceiving() {
	ap.recvLoopOnce.Do(func() {
		ap.log.Debug("starting accesspoint recv loop")
		go ap.recvLoop()

		// set last ping in the future
		ap.lastPongAck = time.Now().Add(pongAckInterval)
		go ap.pongAckTicker()
	})
}

func (ap *Accesspoint) recvLoop() {
loop:
	for {
		select {
		case <-ap.recvLoopStop:
			break loop
		default:
			// no need to hold the connMu since reconnection happens in this routine
			pkt, payload, err := ap.encConn.receivePacket(context.TODO())
			if err != nil {
				if !ap.stop {
					ap.log.Error("failed receiving packet", "error", err)
				}
				break loop
			}

			switch pkt {
			case PacketTypePing:
				if err := ap.encConn.sendPacket(context.TODO(), PacketTypePong, payload); err != nil {
					ap.log.Error("failed sending Pong packet", "error", err)
					break loop
				}
			case PacketTypePongAck:
				ap.lastPongAckLock.Lock()
				ap.lastPongAck = time.Now()
				ap.lastPongAckLock.Unlock()
			default:
				ap.recvChansLock.RLock()
				ll := ap.recvChans[pkt]
				ap.recvChansLock.RUnlock()

				handled := false
				for _, ch := range ll {
					ch <- Packet{Type: pkt, Payload: payload}
					handled = true
				}

				if !handled {
					ap.log.Debug("skipping packet", "type", pkt, "len", len(payload))
				}
			}
		}
	}

	// always close as we might end up here because of application errors
	_ = ap.conn.Close()

	// if we shouldn't stop, try to reconnect
	if !ap.stop {
		ap.connMu.Lock()
		if _, err := backoff.Retry(context.TODO(), func() (struct{}, error) {
			return struct{}{}, ap.reconnect()
		}, backoff.WithBackOff(backoff.NewExponentialBackOff())); err != nil {
			ap.log.Error("failed reconnecting accesspoint", "error", err)
			ap.connMu.Unlock()
			ap.Close()
			return
		}
		ap.connMu.Unlock()
		return
	}

	ap.recvChansLock.RLock()
	defer ap.recvChansLock.RUnlock()

	var closedChannels []chan Packet
	for _, ll := range ap.recvChans {
		for _, ch := range ll {
			if !slices.Contains(closedChannels, ch) {
				closedChannels = append(closedChannels, ch)
				close(ch)
			}
		}
	}
}

func (ap *Accesspoint) pongAckTicker() {
	ticker := time.NewTicker(pongAckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ap.pongAckTickerStop:
			return
		case <-ticker.C:
			ap.lastPongAckLock.Lock()
			timePassed := time.Since(ap.lastPongAck)
			ap.lastPongAckLock.Unlock()
			if timePassed > pongAckInterval {
				ap.log.Error("did not receive last pong ack from accesspoint", "seconds", timePassed.Seconds())
				// closing the connection makes the read in recvLoop fail
				_ = ap.conn.Close()
			}
		}
	}
}

func (ap *Accesspoint) reconnect() error {
	if ap.welcome == nil {
		return backoff.Permanent(fmt.Errorf("cannot reconnect without APWelcome"))
	}

	if err := ap.connect(context.TODO(), pb.LoginCredentials_builder{
		Typ:      ap.welcome.GetReusableAuthCredentialsType().Enum(),
		Username: proto.String(ap.welcome.GetCanonicalUsername()),
		AuthData: ap.welcome.GetReusableAuthCredentials(),
	}.Build()); err != nil {
		return err
	}

	// if we are here the recvLoop has already died, restart it
	go ap.recvLoop()

	ap.log.Debug("re-established accesspoint connection")
	return nil
}

func (ap *Accesspoint) performKeyExchange() ([]byte, error) {
	// accumulate transferred data for challenge
	cc := &connAccumulator{Conn: ap.conn}

	productFlags := []pb.ProductFlags{pb.ProductFlags_PRODUCT_FLAG_NONE}

	// send ClientHello message
	if err := writeMessage(cc, true, pb.ClientHello_builder{
		BuildInfo: pb.BuildInfo_builder{
			Product:      pb.Product_PRODUCT_CLIENT.Enum(),
			ProductFlags: productFlags,
			Platform:     getPlatform().Enum(),
			Version:      proto.Uint64(SpotifyVersionCode),
		}.Build(),
		CryptosuitesSupported: []pb.Cryptosuite{pb.Cryptosuite_CRYPTO_SUITE_SHANNON},
		ClientNonce:           ap.nonce,
		Padding:               []byte{0x1e},
		LoginCryptoHello: pb.LoginCryptoHelloUnion_builder{
			DiffieHellman: pb.LoginCryptoDiffieHellmanHello_builder{
				Gc:              ap.dh.PublicKeyBytes(),
				ServerKeysKnown: proto.Uint32(1),
			}.Build(),
		}.Build(),
	}.Build()); err != nil {
		return nil, fmt.Errorf("failed writing ClientHello message: %w", err)
	}

	// receive APResponseMessage message
	var apResponse pb.APResponseMessage
	if err := readMessage(cc, -1, &apResponse); err != nil {
		return nil, fmt.Errorf("failed reading APResponseMessage message: %w", err)
	}

	dhChallenge := apResponse.GetChallenge().GetLoginCryptoChallenge().GetDiffieHellman()

	// verify signature
	if !verifySignature(dhChallenge.GetGs(), dhChallenge.GetGsSignature()) {
		return nil, fmt.Errorf("failed verifying signature")
	}

	// exchange keys and compute shared secret
	ap.dh.Exchange(dhChallenge.GetGs())

	ap.log.Debug("completed keyexchange")
	return cc.Dump(), nil
}

func (ap *Accesspoint) solveChallenge(exchangeData []byte) error {
	macData := make([]byte, 0, sha1.Size*5)

	mac := hmac.New(sha1.New, ap.dh.SharedSecretBytes())
	for i := byte(1); i < 6; i++ {
		mac.Reset()
		mac.Write(exchangeData)
		mac.Write([]byte{i})
		macData = mac.Sum(macData)
	}

	mac = hmac.New(sha1.New, macData[:20])
	mac.Write(exchangeData)

	if err := writeMessage(ap.conn, false, pb.ClientResponsePlaintext_builder{
		PowResponse:    pb.PoWResponseUnion_builder{}.Build(),
		CryptoResponse: pb.CryptoResponseUnion_builder{}.Build(),
		LoginCryptoResponse: pb.LoginCryptoResponseUnion_builder{
			DiffieHellman: pb.LoginCryptoDiffieHellmanResponse_builder{
				Hmac: mac.Sum(nil),
			}.Build(),
		}.Build(),
	}.Build()); err != nil {
		return fmt.Errorf("failed writing ClientResponsePlaintext message: %w", err)
	}

	// we are not sure if the challenge is actually completed, we check it in authenticate
	ap.encConn = newShannonConn(ap.conn, macData[20:52], macData[52:84])
	ap.log.Debug("completed challenge")
	return nil
}

func (ap *Accesspoint) authenticate(ctx context.Context, credentials *pb.LoginCredentials) error {
	if ap.encConn == nil {
		panic("accesspoint not connected")
	}

	// assemble ClientResponseEncrypted message
	payload, err := proto.Marshal(pb.ClientResponseEncrypted_builder{
		LoginCredentials: credentials,
		VersionString:    proto.String(versionString()),
		SystemInfo: pb.SystemInfo_builder{
			Os:                      getOS().Enum(),
			CpuFamily:               getCpuFamily().Enum(),
			SystemInformationString: proto.String(systemInfoString()),
			DeviceId:                proto.String(ap.deviceId),
		}.Build(),
	}.Build())
	if err != nil {
		return fmt.Errorf("failed marshalling ClientResponseEncrypted message: %w", err)
	}

	// send Login packet
	if err := ap.encConn.sendPacket(ctx, PacketTypeLogin, payload); err != nil {
		return fmt.Errorf("failed sending Login packet: %w", err)
	}

	// check if we received an APResponseMessage from the challenge
	var challengeResp pb.APResponseMessage
	if peekBytes, err := ap.encConn.peekUnencrypted(9); err != nil {
		return fmt.Errorf("failed peeking unencrypted bytes: %w", err)
	} else if err = readMessage(bytes.NewReader(peekBytes), 9, &challengeResp); err == nil {
		return &AccesspointLoginError{Message: challengeResp.GetLoginFailed()}
	}

	// receive APWelcome or AuthFailure
	recvPkt, recvPayload, err := ap.encConn.receivePacket(ctx)
	if err != nil {
		return fmt.Errorf("failed receiving Login response packet: %w", err)
	}

	switch recvPkt {
	case PacketTypeAPWelcome:
		var welcome pb.APWelcome
		if err := proto.Unmarshal(recvPayload, &welcome); err != nil {
			return fmt.Errorf("failed unmarshalling APWelcome message: %w", err)
		}

		ap.welcome = &welcome
		ap.log.Info("authenticated AP", "username", obfuscateUsername(welcome.GetCanonicalUsername()))
		return nil
	case PacketTypeAuthFailure:
		var loginFailed pb.APLoginFailed
		if err := proto.Unmarshal(recvPayload, &loginFailed); err != nil {
			return fmt.Errorf("failed unmarshalling APLoginFailed message: %w", err)
		}
		return &AccesspointLoginError{Message: &loginFailed}
	default:
		return fmt.Errorf("unexpected command after Login packet: %x", recvPkt)
	}
}
