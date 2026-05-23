package connect

// connectBody is a sealed JSON request body for connect APIs.
type connectBody interface {
	iConnectBody()
}

// endpointPayload is a sealed JSON body for connect-state player commands.
type endpointPayload interface {
	connectBody
	iEndpoint()
}

// LoggingParams mirrors connect-state command logging_params.
type LoggingParams struct {
	CommandID string `json:"command_id"`
}

func newLoggingParams() LoggingParams {
	return LoggingParams{CommandID: randomID()}
}

type playerCommandHeader struct {
	Endpoint      Endpoint      `json:"endpoint"`
	LoggingParams LoggingParams `json:"logging_params"`
}

// SimpleCommandPayload is the body for endpoint-only player commands.
type SimpleCommandPayload struct {
	Command SimplePlayerCommand `json:"command"`
}

type SimplePlayerCommand struct {
	playerCommandHeader
}

func (SimpleCommandPayload) iConnectBody() {}
func (SimpleCommandPayload) iEndpoint()    {}

func newSimpleCommand(endpoint Endpoint) SimpleCommandPayload {
	return SimpleCommandPayload{
		Command: SimplePlayerCommand{
			playerCommandHeader: playerCommandHeader{
				Endpoint:      endpoint,
				LoggingParams: newLoggingParams(),
			},
		},
	}
}

// SeekCommandPayload is the body for seek_to.
type SeekCommandPayload struct {
	Command SeekPlayerCommand `json:"command"`
}

type SeekPlayerCommand struct {
	playerCommandHeader
	Value int `json:"value"`
}

func (SeekCommandPayload) iConnectBody() {}
func (SeekCommandPayload) iEndpoint()    {}

func newSeekCommand(positionMS int) SeekCommandPayload {
	return SeekCommandPayload{
		Command: SeekPlayerCommand{
			playerCommandHeader: playerCommandHeader{
				Endpoint:      EndpointSeekTo,
				LoggingParams: newLoggingParams(),
			},
			Value: positionMS,
		},
	}
}

// ShuffleCommandPayload is the body for set_shuffling_context.
type ShuffleCommandPayload struct {
	Command ShufflePlayerCommand `json:"command"`
}

type ShufflePlayerCommand struct {
	playerCommandHeader
	Value bool `json:"value"`
}

func (ShuffleCommandPayload) iConnectBody() {}
func (ShuffleCommandPayload) iEndpoint()    {}

func newShuffleCommand(enabled bool) ShuffleCommandPayload {
	return ShuffleCommandPayload{
		Command: ShufflePlayerCommand{
			playerCommandHeader: playerCommandHeader{
				Endpoint:      EndpointSetShufflingContext,
				LoggingParams: newLoggingParams(),
			},
			Value: enabled,
		},
	}
}

// SetOptionsCommandPayload is the body for set_options (repeat mode).
type SetOptionsCommandPayload struct {
	Command SetOptionsPlayerCommand `json:"command"`
}

type SetOptionsPlayerCommand struct {
	playerCommandHeader
	RepeatingTrack   bool `json:"repeating_track"`
	RepeatingContext bool `json:"repeating_context"`
}

func (SetOptionsCommandPayload) iConnectBody() {}
func (SetOptionsCommandPayload) iEndpoint()    {}

func newSetOptionsCommand(repeatingTrack, repeatingContext bool) SetOptionsCommandPayload {
	return SetOptionsCommandPayload{
		Command: SetOptionsPlayerCommand{
			playerCommandHeader: playerCommandHeader{
				Endpoint:      EndpointSetOptions,
				LoggingParams: newLoggingParams(),
			},
			RepeatingTrack:   repeatingTrack,
			RepeatingContext: repeatingContext,
		},
	}
}

// AddToQueueCommandPayload is the body for add_to_queue.
type AddToQueueCommandPayload struct {
	Command AddToQueuePlayerCommand `json:"command"`
}

type QueueTrackRef struct {
	URI string `json:"uri"`
}

type AddToQueuePlayerCommand struct {
	playerCommandHeader
	Track QueueTrackRef `json:"track"`
}

func (AddToQueueCommandPayload) iConnectBody() {}
func (AddToQueueCommandPayload) iEndpoint()    {}

func newAddToQueueCommand(uri string) AddToQueueCommandPayload {
	return AddToQueueCommandPayload{
		Command: AddToQueuePlayerCommand{
			playerCommandHeader: playerCommandHeader{
				Endpoint:      EndpointAddToQueue,
				LoggingParams: newLoggingParams(),
			},
			Track: QueueTrackRef{URI: uri},
		},
	}
}

// PlayCommandPayload is the body for play with a context or track uri.
type PlayCommandPayload struct {
	Command PlayPlayerCommand `json:"command"`
}

type PlayContextRef struct {
	URI string `json:"uri"`
	URL string `json:"url"`
}

type PlaySkipTo struct {
	TrackURI string `json:"track_uri"`
}

type PlayCommandOptions struct {
	SkipTo PlaySkipTo `json:"skip_to"`
}

type PlayPlayerCommand struct {
	playerCommandHeader
	Context PlayContextRef      `json:"context"`
	Options *PlayCommandOptions `json:"options,omitempty"`
}

func (PlayCommandPayload) iConnectBody() {}
func (PlayCommandPayload) iEndpoint()    {}

func newPlayCommand(uri string) PlayCommandPayload {
	cmd := PlayPlayerCommand{
		playerCommandHeader: playerCommandHeader{
			Endpoint:      EndpointPlay,
			LoggingParams: newLoggingParams(),
		},
		Context: PlayContextRef{
			URI: uri,
			URL: "context://" + uri,
		},
	}
	if !IsContextURI(uri) {
		cmd.Options = &PlayCommandOptions{
			SkipTo: PlaySkipTo{TrackURI: uri},
		}
	}
	return PlayCommandPayload{Command: cmd}
}

// VolumePayload is the body for PUT .../connect/volume/...
type VolumePayload struct {
	Volume int `json:"volume"`
}

func (VolumePayload) iConnectBody() {}

// TransferPayload is the body for POST .../connect/transfer/...
type TransferPayload struct {
	TransferOptions TransferOptions `json:"transfer_options"`
	CommandID       string          `json:"command_id"`
}

type TransferRestore string

const TransferRestoreResume TransferRestore = "resume"

type TransferOptions struct {
	RestorePaused TransferRestore `json:"restore_paused"`
}

func (TransferPayload) iConnectBody() {}

func newTransferPayload() TransferPayload {
	return TransferPayload{
		TransferOptions: TransferOptions{RestorePaused: TransferRestoreResume},
		CommandID:       randomID(),
	}
}

// StateRequestPayload is the body for PUT .../devices/hobs_* (connect-state poll).
type StateRequestPayload struct {
	MemberType StateRequestMemberType `json:"member_type"`
	Device     StateRequestDevice     `json:"device"`
}

type StateRequestDevice struct {
	DeviceInfo StateRequestDeviceInfo `json:"device_info"`
}

type StateRequestDeviceInfo struct {
	Capabilities StateRequestCapabilities `json:"capabilities"`
}

type StateRequestCapabilities struct {
	CanBePlayer          bool `json:"can_be_player"`
	Hidden               bool `json:"hidden"`
	NeedsFullPlayerState bool `json:"needs_full_player_state"`
}

type StateRequestMemberType string

const StateRequestMemberTypeConnectState StateRequestMemberType = "CONNECT_STATE"

func (StateRequestPayload) iConnectBody() {}

func newStateRequestPayload() StateRequestPayload {
	return StateRequestPayload{
		MemberType: StateRequestMemberTypeConnectState,
		Device: StateRequestDevice{
			DeviceInfo: StateRequestDeviceInfo{
				Capabilities: StateRequestCapabilities{
					CanBePlayer:          false,
					Hidden:               true,
					NeedsFullPlayerState: true,
				},
			},
		},
	}
}

// Register device wire constants.
const (
	DeviceTypeComputer             = "computer"
	BrandSpotify                   = "spotify"
	PlayTokenLostPause             = "pause"
	ManifestFormatFileIDsMP3       = "file_ids_mp3"
	ManifestFormatFileURLsMP3      = "file_urls_mp3"
	ManifestFormatFileIDsMP4       = "file_ids_mp4"
	ManifestFormatManifestIDsVideo = "manifest_ids_video"
)

// RegisterDevicePayload is the body for POST .../track-playback/v1/devices.
type RegisterDevicePayload struct {
	Device                  RegisterDeviceBody `json:"device"`
	OutroEndcontentSnooping bool               `json:"outro_endcontent_snooping"`
	ConnectionID            string             `json:"connection_id"`
	ClientVersion           string             `json:"client_version"`
	Volume                  int                `json:"volume"`
}

type RegisterDeviceBody struct {
	DeviceID           string                     `json:"device_id"`
	DeviceType         string                     `json:"device_type"`
	Brand              string                     `json:"brand"`
	Model              string                     `json:"model"`
	Name               string                     `json:"name"`
	IsGroup            bool                       `json:"is_group"`
	Metadata           struct{}                   `json:"metadata"`
	PlatformIdentifier string                     `json:"platform_identifier"`
	Capabilities       RegisterDeviceCapabilities `json:"capabilities"`
}

type RegisterDeviceCapabilities struct {
	ChangeVolume          bool     `json:"change_volume"`
	SupportsFileMediaType bool     `json:"supports_file_media_type"`
	EnablePlayToken       bool     `json:"enable_play_token"`
	PlayTokenLostBehavior string   `json:"play_token_lost_behavior"`
	DisableConnect        bool     `json:"disable_connect"`
	AudioPodcasts         bool     `json:"audio_podcasts"`
	VideoPlayback         bool     `json:"video_playback"`
	ManifestFormats       []string `json:"manifest_formats"`
}

func (RegisterDevicePayload) iConnectBody() {}

func newRegisterDevicePayload(deviceID, connectionID, platformID string) RegisterDevicePayload {
	return RegisterDevicePayload{
		Device: RegisterDeviceBody{
			DeviceID:           deviceID,
			DeviceType:         DeviceTypeComputer,
			Brand:              BrandSpotify,
			Model:              connectDeviceModel,
			Name:               connectDeviceName,
			IsGroup:            false,
			PlatformIdentifier: platformID,
			Capabilities: RegisterDeviceCapabilities{
				ChangeVolume:          true,
				SupportsFileMediaType: true,
				EnablePlayToken:       true,
				PlayTokenLostBehavior: PlayTokenLostPause,
				DisableConnect:        false,
				AudioPodcasts:         true,
				VideoPlayback:         true,
				ManifestFormats: []string{
					ManifestFormatFileIDsMP3,
					ManifestFormatFileURLsMP3,
					ManifestFormatFileIDsMP4,
					ManifestFormatManifestIDsVideo,
				},
			},
		},
		OutroEndcontentSnooping: false,
		ConnectionID:            connectionID,
		ClientVersion:           clientVersion,
		Volume:                  65535,
	}
}
