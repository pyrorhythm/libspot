package connect

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// State returns the current connect-state snapshot.
func (c *Connect) State(ctx context.Context) (State, error) {
	return c.fetchState(ctx)
}

// Playback returns the current Connect player state.
func (c *Connect) Playback(ctx context.Context) (Playback, error) {
	return withState(ctx, c, func(state State) (Playback, error) {
		return mapPlayback(state), nil
	})
}

// Devices lists the devices visible in the current Connect session.
func (c *Connect) Devices(ctx context.Context) ([]Device, error) {
	return withState(ctx, c, func(state State) ([]Device, error) {
		return state.Devices, nil
	})
}

// Queue returns the currently playing item plus the upcoming tracks.
func (c *Connect) Queue(ctx context.Context) (Queue, error) {
	return withState(ctx, c, func(state State) (Queue, error) {
		return mapQueue(state), nil
	})
}

// Play resumes playback, or starts the given uri when one is supplied.
func (c *Connect) Play(ctx context.Context, uri string) error {
	return withStateErr(ctx, c, func(state State) error {
		if state.ActiveDeviceID == "" {
			targetID := resolveTargetDeviceID(state, c.opts.Device)
			if targetID == "" {
				return errMissingDevice
			}
			state.ActiveDeviceID = targetID
		}
		if uri == "" {
			return c.sendPlayerCommand(ctx, state, "resume", nil)
		}
		return c.sendPlayerCommand(ctx, state, "play", newPlayCommand(uri))
	})
}

func (c *Connect) Pause(ctx context.Context) error {
	return c.sendDirectCommand(ctx, "pause", nil)
}

func (c *Connect) Next(ctx context.Context) error {
	return c.sendDirectCommand(ctx, "skip_next", nil)
}

func (c *Connect) Previous(ctx context.Context) error {
	return c.sendDirectCommand(ctx, "skip_prev", nil)
}

func (c *Connect) Seek(ctx context.Context, positionMS int) error {
	if positionMS < 0 {
		positionMS = 0
	}
	return c.sendDirectCommand(ctx, "seek_to", newSeekCommand(positionMS))
}

func (c *Connect) Volume(ctx context.Context, volume int) error {
	volume = clampVolume(volume)
	return withStateErr(ctx, c, func(state State) error {
		fromID := transferSourceID(state)
		if fromID == "" || state.ActiveDeviceID == "" {
			return errMissingDevice
		}
		base, err := c.connectStateBase()
		if err != nil {
			return err
		}
		url := fmt.Sprintf(
			"%s/connect/volume/from/%s/to/%s",
			base,
			fromID,
			state.ActiveDeviceID,
		)
		return c.sendConnectRequest(ctx, http.MethodPut, url, VolumePayload{
			Volume: int(float64(volume) / 100 * 65535),
		})
	})
}

func (c *Connect) Shuffle(ctx context.Context, enabled bool) error {
	return c.sendDirectCommand(ctx, "set_shuffling_context", newShuffleCommand(enabled))
}

// Repeat sets the repeat mode: "track", "context" or anything else for off.
func (c *Connect) Repeat(ctx context.Context, mode string) error {
	repeatingTrack, repeatingContext := repeatFlags(mode)
	return c.sendDirectCommand(ctx, "set_options", newSetOptionsCommand(repeatingTrack, repeatingContext))
}

func (c *Connect) QueueAdd(ctx context.Context, uri string) error {
	return c.sendDirectCommand(ctx, "add_to_queue", newAddToQueueCommand(uri))
}

// Transfer moves playback to the device identified by deviceID.
func (c *Connect) Transfer(ctx context.Context, deviceID string) error {
	return withStateErr(ctx, c, func(state State) error {
		fromID := transferSourceID(state)
		if fromID == "" {
			return errMissingDevice
		}
		base, err := c.connectStateBase()
		if err != nil {
			return err
		}
		url := fmt.Sprintf(
			"%s/connect/transfer/from/%s/to/%s",
			base,
			fromID,
			deviceID,
		)
		return c.sendConnectCommand(ctx, url, newTransferPayload())
	})
}

func (c *Connect) sendDirectCommand(
	ctx context.Context,
	endpoint string,
	payload any,
) error {
	if payload == nil {
		payload = newSimpleCommand(endpoint)
	}
	fromID, toID, ok := c.commandRoute()
	if !ok {
		return c.sendStateCommand(ctx, endpoint, payload)
	}
	base, err := c.connectStateBase()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/player/command/from/%s/to/%s", base, fromID, toID)
	if err := c.sendConnectCommand(ctx, url, payload); err != nil {
		if !isRouteStaleError(err) {
			return err
		}
		return c.sendStateCommand(ctx, endpoint, payload)
	}
	return nil
}

func (c *Connect) sendStateCommand(
	ctx context.Context,
	endpoint string,
	payload any,
) error {
	return withStateErr(ctx, c, func(state State) error {
		return c.sendPlayerCommand(ctx, state, endpoint, payload)
	})
}

func withState[T any](
	ctx context.Context,
	c *Connect,
	fn func(State) (T, error),
) (T, error) {
	state, err := c.fetchState(ctx)
	if err != nil {
		var zero T
		return zero, err
	}
	return fn(state)
}

func withStateErr(ctx context.Context, c *Connect, fn func(State) error) error {
	_, err := withState(ctx, c, func(state State) (struct{}, error) {
		return struct{}{}, fn(state)
	})
	return err
}

func transferSourceID(state State) string {
	if state.OriginDeviceID != "" {
		return state.OriginDeviceID
	}
	return state.ActiveDeviceID
}

func resolveTargetDeviceID(state State, selector string) string {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return ""
	}
	for _, device := range state.Devices {
		if strings.EqualFold(device.ID, selector) || strings.EqualFold(device.Name, selector) {
			return device.ID
		}
	}
	return ""
}

func clampVolume(volume int) int {
	if volume < 0 {
		return 0
	}
	if volume > 100 {
		return 100
	}
	return volume
}

func repeatFlags(mode string) (track, context bool) {
	switch strings.ToLower(mode) {
	case "track":
		return true, false
	case "context":
		return false, true
	default:
		return false, false
	}
}
