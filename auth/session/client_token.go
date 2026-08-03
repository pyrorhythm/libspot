package session

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"google.golang.org/protobuf/proto"
	"pyrorhythm.dev/libspot"
	datav0 "pyrorhythm.dev/libspot/gen/spotify/clienttoken/data/v0"
	httpv0 "pyrorhythm.dev/libspot/gen/spotify/clienttoken/http/v0"
	"pyrorhythm.dev/libspot/pkg/transport"
	"resty.dev/v3"
)

const ctUrl = "https://clienttoken.spotify.com/v1/clienttoken"

// retrieveClientToken fetches a Spotify client token using the client-token API.
//
// Sourced from devgianlu/go-librespot.
func retrieveClientToken(deviceId string) (string, error) {
	body, err := proto.Marshal(clientTokenRequest(deviceId))
	if err != nil {
		return "", fmt.Errorf("failed marshalling ClientTokenRequest: %w", err)
	}

	resp, err := resty.NewWithClient(transport.HTTPClient(0)).R().
		SetHeaderMultiValues(map[string][]string{
			"Accept":     {"application/x-protobuf"},
			"User-Agent": {fmt.Sprintf("libspot/0.0.0 Go/%s", runtime.Version())},
		}).SetBody(body).Post(ctUrl)
	if err != nil {
		return "", fmt.Errorf("failed requesting clienttoken: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("invalid status code from clienttoken: %d", resp.StatusCode())
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
			ClientId:      libspot.ClientIdHex,
			ClientVersion: "0.0.0",
			ConnectivitySdkData: datav0.ConnectivitySdkData_builder{
				DeviceId:             deviceID,
				PlatformSpecificData: platformSpecificData(),
			}.Build(),
		}.Build(),
	}.Build()
}

func platformSpecificData() *datav0.PlatformSpecificData {
	psd := datav0.PlatformSpecificData_builder{}

	switch runtime.GOOS {
	case "android":
		psd.Android = &datav0.NativeAndroidData{}
	case "darwin":
		psd.Mac = &datav0.NativeDesktopMacOSData{}
	case "ios":
		psd.Ios = &datav0.NativeIOSData{}
	case "linux", "freebsd":
		psd.Linux = &datav0.NativeDesktopLinuxData{}
	case "windows":
		psd.Win = &datav0.NativeDesktopWindowsData{}
	}

	return psd.Build()
}
