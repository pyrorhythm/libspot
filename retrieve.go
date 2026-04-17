package libspot

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	datav0 "github.com/pyrorhythm/libspot/gen/spotify/clienttoken/data/v0"
	httpv0 "github.com/pyrorhythm/libspot/gen/spotify/clienttoken/http/v0"
	"google.golang.org/protobuf/proto"
	"resty.dev/v3"
)

const ctUrl = "https://clienttoken.spotify.com/v1/clienttoken"

// RetrieveClientToken fetches a Spotify client token using the client-token API.
//
// Sourced from devgianlu/go-librespot.
func RetrieveClientToken(deviceId string) (string, error) {
	body, err := proto.Marshal(clientTokenRequest(deviceId))
	if err != nil {
		return "", fmt.Errorf("failed marshalling ClientTokenRequest: %w", err)
	}

	resp, err := resty.New().R().
		SetHeaderMultiValues(map[string][]string{
			"Accept":     {"application/x-protobuf"},
			"User-Agent": {fmt.Sprintf("libspot/0.0.0 Go/%s", runtime.Version())},
		}).SetBody(body).Post(ctUrl)
	if err != nil {
		return "", fmt.Errorf("failed requesting clienttoken: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("invalid status code from clienttoken: %d", resp.StatusCode)
	}

	var protoResp httpv0.ClientTokenResponse
	if err := proto.Unmarshal(resp.Bytes(), &protoResp); err != nil {
		return "", fmt.Errorf("failed unmarshalling clienttoken response: %w", err)
	}

	switch protoResp.GetResponseType() {
	case httpv0.ClientTokenResponseType_RESPONSE_GRANTED_TOKEN_RESPONSE:
		granted := protoResp.GetGrantedToken()
		if granted == nil {
			return "", errors.New("invalid granted token response")
		}
		return granted.GetToken(), nil
	case httpv0.ClientTokenResponseType_RESPONSE_CHALLENGES_RESPONSE:
		return "", errors.New("clienttoken challenge not supported")
	default:
		return "", fmt.Errorf("unknown clienttoken response type: %v", protoResp.GetResponseType())
	}
}

func clientTokenRequest(deviceID string) *httpv0.ClientTokenRequest {
	return httpv0.ClientTokenRequest_builder{
		RequestType: httpv0.ClientTokenRequestType_REQUEST_CLIENT_DATA_REQUEST,
		ClientData: httpv0.ClientDataRequest_builder{
			ClientId:      ClientIdHex,
			ClientVersion: "0.0.0",
			ConnectivitySdkData: datav0.ConnectivitySdkData_builder{
				DeviceId:             deviceID,
				PlatformSpecificData: platformSpecificData(),
			}.Build(),
		}.Build(),
	}.Build()
}
