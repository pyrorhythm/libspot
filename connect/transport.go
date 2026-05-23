package connect

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// fetchState registers our hidden member (if needed) and PUTs a full-state
// request, returning the decoded public snapshot.
func (c *Connect) fetchState(ctx context.Context) (State, error) {
	connectionID, err := c.ensureConnectDevice(ctx)
	if err != nil {
		return State{}, err
	}
	c.mu.Lock()
	deviceID := c.connectDeviceID
	c.mu.Unlock()

	base, err := c.connectStateBase()
	if err != nil {
		return State{}, err
	}

	payload := newStateRequestPayload()

	resp, err := c.http.R().SetContext(ctx).
		SetMethod(http.MethodPut).
		SetURL(fmt.Sprintf("%s/devices/hobs_%s", base, deviceID)).
		SetHeaders(connectHeaders()).
		SetHeader("x-spotify-connection-id", connectionID).
		SetBody(mustJSON(payload)).
		Send()
	if err != nil {
		return State{}, errors.Wrap(err, "connect: fetch state")
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return State{}, newAPIError(resp.StatusCode(), resp.Bytes())
	}

	state, err := parseState(resp.Bytes())
	if err != nil {
		return state, err
	}
	c.cacheCommandRoute(state)
	return state, nil
}

func (c *Connect) ensureConnectDevice(ctx context.Context) (string, error) {
	connectionID, err := c.connectionID()
	if err != nil {
		return "", err
	}

	c.mu.Lock()
	if c.connectDeviceID == "" {
		c.connectDeviceID = randomID()
	}
	needs := c.lastCID != connectionID || time.Since(c.registeredAt) > connectionTTL
	c.mu.Unlock()
	if !needs {
		return connectionID, nil
	}

	if err := c.registerDevice(ctx, connectionID); err != nil {
		return "", err
	}

	c.mu.Lock()
	c.lastCID = connectionID
	c.registeredAt = time.Now()
	c.mu.Unlock()
	return connectionID, nil
}

func (c *Connect) registerDevice(ctx context.Context, connectionID string) error {
	c.mu.Lock()
	deviceID := c.connectDeviceID
	c.mu.Unlock()

	base, err := c.trackPlaybackBase()
	if err != nil {
		return err
	}

	platformID := fmt.Sprintf("web_player %s;libspot", runtime.GOOS)
	payload := newRegisterDevicePayload(deviceID, connectionID, platformID)

	resp, err := c.http.R().SetContext(ctx).
		SetMethod(http.MethodPost).
		SetURL(base + "/devices").
		SetHeaders(connectHeaders()).
		SetBody(mustJSON(payload)).
		Send()
	if err != nil {
		c.invalidateCommandRoute()
		return errors.Wrap(err, "connect: register device")
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		c.invalidateCommandRoute()
		return newAPIError(resp.StatusCode(), resp.Bytes())
	}
	return nil
}

func (c *Connect) ownedDevice() (deviceID, connectionID string, ok bool) {
	c.mu.Lock()
	deviceID = c.connectDeviceID
	c.mu.Unlock()
	if deviceID == "" {
		return "", "", false
	}
	// Best-effort during shutdown; connect-state DELETE still works without it.
	connectionID, _ = c.connectionID()
	return deviceID, connectionID, true
}

func connectStateDeviceHeaders(connectionID string) map[string]string {
	headers := connectHeaders()
	if connectionID != "" {
		headers["x-spotify-connection-id"] = connectionID
	}
	return headers
}

func hobsDeviceID(deviceID string) string {
	return fmt.Sprintf("hobs_%s", deviceID)
}

func (c *Connect) markInactive(
	ctx context.Context,
	deviceID, connectionID string,
	notify bool,
) error {
	base, err := c.connectStateBase()
	if err != nil {
		return err
	}
	url := fmt.Sprintf(
		"%s/devices/%s/inactive?notify=%s",
		base,
		hobsDeviceID(deviceID),
		strconv.FormatBool(notify),
	)
	err = c.doRequest(ctx, http.MethodPut, url, connectStateDeviceHeaders(connectionID), nil)
	if err != nil && !isBenignDeviceError(err) {
		return errors.Wrap(err, "connect: mark inactive")
	}
	return nil
}

func (c *Connect) deleteConnectStateDevice(
	ctx context.Context,
	deviceID, connectionID string,
) error {
	base, err := c.connectStateBase()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/devices/%s", base, hobsDeviceID(deviceID))
	err = c.doRequest(ctx, http.MethodDelete, url, connectStateDeviceHeaders(connectionID), nil)
	if err != nil && !isBenignDeviceError(err) {
		return errors.Wrap(err, "connect: delete connect-state device")
	}
	return nil
}

func (c *Connect) deleteTrackPlaybackDevice(ctx context.Context, deviceID string) error {
	base, err := c.trackPlaybackBase()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/devices/%s", base, deviceID)
	err = c.doRequest(ctx, http.MethodDelete, url, connectHeaders(), nil)
	if err != nil && !isBenignDeviceError(err) {
		return errors.Wrap(err, "connect: delete track-playback device")
	}
	return nil
}

func (c *Connect) forgetLocalDevice() {
	c.mu.Lock()
	c.connectDeviceID = ""
	c.lastCID = ""
	c.registeredAt = time.Time{}
	c.mu.Unlock()
	c.invalidateCommandRoute()
}

func (c *Connect) doRequest(
	ctx context.Context,
	method, url string,
	headers map[string]string,
	body []byte,
) error {
	req := c.http.R().SetContext(ctx).SetMethod(method).SetURL(url).SetHeaders(headers)
	if body != nil {
		req = req.SetBody(body)
	}
	resp, err := req.Send()
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return newAPIError(resp.StatusCode(), resp.Bytes())
}

func (c *Connect) sendPlayerCommand(
	ctx context.Context,
	state State,
	payload endpointPayload,
) error {
	fromID := state.OriginDeviceID
	if fromID == "" {
		c.mu.Lock()
		fromID = c.connectDeviceID
		c.mu.Unlock()
	}
	if fromID == "" {
		fromID = state.ActiveDeviceID
	}
	if fromID == "" || state.ActiveDeviceID == "" {
		return errMissingDevice
	}
	base, err := c.connectStateBase()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/player/command/from/%s/to/%s", base, fromID, state.ActiveDeviceID)
	return c.sendConnectRequest(ctx, http.MethodPost, url, payload)
}

func (c *Connect) sendConnectCommand(
	ctx context.Context,
	url string,
	payload connectBody,
) error {
	return c.sendConnectRequest(ctx, http.MethodPost, url, payload)
}

func (c *Connect) sendConnectRequest(
	ctx context.Context,
	method, url string,
	payload connectBody,
) error {
	err := c.doRequest(ctx, method, url, connectHeaders(), mustJSON(payload))
	if err != nil {
		c.invalidateCommandRoute()
		return errors.Wrap(err, "connect: send command")
	}
	return nil
}
