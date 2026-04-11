package libspot

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"

	datav0 "github.com/pyrorhythm/libspot/api/spotify/clienttoken/data/v0"
	httpv0 "github.com/pyrorhythm/libspot/api/spotify/clienttoken/http/v0"
	"google.golang.org/protobuf/proto"
)

// RetrieveClientToken fetches a Spotify client token using the client-token API.
func RetrieveClientToken(c *http.Client, deviceId string) (string, error) {
	body, err := proto.Marshal(clientTokenRequest(deviceId))
	if err != nil {
		return "", fmt.Errorf("failed marshalling ClientTokenRequest: %w", err)
	}

	reqURL, err := url.Parse("https://clienttoken.spotify.com/v1/clienttoken")
	if err != nil {
		return "", fmt.Errorf("invalid clienttoken url: %w", err)
	}

	resp, err := c.Do(&http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: http.Header{
			"Accept":     []string{"application/x-protobuf"},
			"User-Agent": []string{fmt.Sprintf("libspot/0.0.0 Go/%s", runtime.Version())},
		},
		Body: io.NopCloser(bytes.NewReader(body)),
	})
	if err != nil {
		return "", fmt.Errorf("failed requesting clienttoken: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code from clienttoken: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading clienttoken response: %w", err)
	}

	var protoResp httpv0.ClientTokenResponse
	if err := proto.Unmarshal(respBody, &protoResp); err != nil {
		return "", fmt.Errorf("failed unmarshalling clienttoken response: %w", err)
	}

	switch protoResp.ResponseType {
	case httpv0.ClientTokenResponseType_RESPONSE_GRANTED_TOKEN_RESPONSE:
		granted := protoResp.GetGrantedToken()
		if granted == nil {
			return "", errors.New("invalid granted token response")
		}
		return granted.Token, nil
	case httpv0.ClientTokenResponseType_RESPONSE_CHALLENGES_RESPONSE:
		return "", errors.New("clienttoken challenge not supported")
	default:
		return "", fmt.Errorf("unknown clienttoken response type: %v", protoResp.ResponseType)
	}
}

func clientTokenRequest(deviceID string) *httpv0.ClientTokenRequest {
	return &httpv0.ClientTokenRequest{
		RequestType: httpv0.ClientTokenRequestType_REQUEST_CLIENT_DATA_REQUEST,
		Request: &httpv0.ClientTokenRequest_ClientData{
			ClientData: &httpv0.ClientDataRequest{
				ClientId:      ClientIdHex,
				ClientVersion: "0.0.0",
				Data: &httpv0.ClientDataRequest_ConnectivitySdkData{
					ConnectivitySdkData: &datav0.ConnectivitySdkData{
						DeviceId:             deviceID,
						PlatformSpecificData: platformSpecificData(),
					},
				},
			},
		},
	}
}
