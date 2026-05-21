package connect

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
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

func (c *Connect) sendPlayerCommand(
	ctx context.Context,
	state State,
	endpoint string,
	payload any,
) error {
	if payload == nil {
		payload = newSimpleCommand(endpoint)
	}
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
	payload any,
) error {
	return c.sendConnectRequest(ctx, http.MethodPost, url, payload)
}

func (c *Connect) sendConnectRequest(
	ctx context.Context,
	method, url string,
	payload any,
) error {
	resp, err := c.http.R().SetContext(ctx).
		SetMethod(method).
		SetURL(url).
		SetHeaders(connectHeaders()).
		SetBody(mustJSON(payload)).
		Send()
	if err != nil {
		c.invalidateCommandRoute()
		return errors.Wrap(err, "connect: send command")
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		c.invalidateCommandRoute()
		return newAPIError(resp.StatusCode(), resp.Bytes())
	}
	return nil
}
