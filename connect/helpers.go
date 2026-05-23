package connect

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

func randomID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}

func mustJSON(payload connectBody) []byte {
	data, _ := json.Marshal(payload)
	return data
}

// isBenignDeviceError reports API responses that mean the device or route is already gone.
func isBenignDeviceError(err error) bool {
	var e APIError
	if !errors.As(err, &e) {
		return false
	}
	switch e.Status {
	case http.StatusBadRequest, http.StatusNotFound, http.StatusConflict, http.StatusGone:
		return true
	default:
		return false
	}
}

func (c *Connect) cacheCommandRoute(state State) {
	c.routeMu.Lock()
	defer c.routeMu.Unlock()
	c.activeID = state.ActiveDeviceID
	c.originID = state.OriginDeviceID
	c.routeStored = state.ActiveDeviceID != ""
}

func (c *Connect) invalidateCommandRoute() {
	c.routeMu.Lock()
	defer c.routeMu.Unlock()
	c.routeStored = false
}

func (c *Connect) commandRoute() (from, to string, ok bool) {
	c.routeMu.RLock()
	active, origin, stored := c.activeID, c.originID, c.routeStored
	c.routeMu.RUnlock()
	if !stored || active == "" {
		return "", "", false
	}
	from = origin
	if from == "" {
		c.mu.Lock()
		from = c.connectDeviceID
		c.mu.Unlock()
	}
	if from == "" {
		from = active
	}
	if from == "" {
		return "", "", false
	}
	return from, active, true
}
