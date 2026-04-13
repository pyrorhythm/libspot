package dealer

import (
	"github.com/pyrorhythm/libspot/dealer/types"
	"github.com/pyrorhythm/libspot/gen/spotify/connectstate"
	"github.com/pyrorhythm/libspot/gen/spotify/parental_controls"
	"github.com/pyrorhythm/libspot/gen/spotify/playbacksettings/pubsub"
	"github.com/pyrorhythm/libspot/gen/spotify/playlist4"
)

// TypedDecoder converts a raw dealer Message into a concrete typed value T.
// Custom topics should use the exported helpers (DecodePB, DecodeJSON,
// DecodeBytes) when building their own decoder rather than hand-rolling
// base64/gzip/proto handling.
type TypedDecoder[T any] func(m *types.Message) (T, error)

// Topic is the single source of truth binding a dealer URI prefix to the
// Go type delivered on it. Topics are the sole public API for typed
// subscription; creating a custom topic is as simple as declaring a
// Topic[T] var with a URI and a decoder.
type Topic[T any] struct {
	URI    string
	Decode TypedDecoder[T]
}

// Subscribe delivers decoded values from topic to cb.
// Decode errors are dropped silently.
func Subscribe[T any](d *Dealer, topic Topic[T], cb func(T)) (unsubscribe func()) {
	return d.OnMsg(topic.URI, func(m *types.Message) {
		v, err := topic.Decode(m)
		if err != nil {
			return
		}
		cb(v)
	})
}

var (
	TopicConnectionID = Topic[string]{
		URI:    connectionIDURIPrefix,
		Decode: decodeConnectionID,
	}
	TopicClusterUpdate = Topic[*connectstate.ClusterUpdate]{
		URI:    "hm://connect-state/v1/cluster",
		Decode: DecodePB[*connectstate.ClusterUpdate],
	}
	TopicSetVolume = Topic[*connectstate.SetVolumeCommand]{
		URI:    "hm://connect-state/v1/connect/volume",
		Decode: DecodePB[*connectstate.SetVolumeCommand],
	}
	TopicLogout = Topic[*connectstate.LogoutCommand]{
		URI:    "hm://connect-state/v1/connect/logout",
		Decode: DecodePB[*connectstate.LogoutCommand],
	}
	TopicPlaylistModification = Topic[*playlist4.PlaylistModificationInfo]{
		URI:    "hm://playlist/v2/playlist/",
		Decode: DecodePB[*playlist4.PlaylistModificationInfo],
	}
	TopicSessionUpdate = Topic[*types.SessionUpdate]{
		URI:    "social-connect/v2/session_update",
		Decode: DecodeJSON[types.SessionUpdate],
	}
	TopicUserAttributesUpdate = Topic[*parental_controls.UserAttributesUpdate]{
		URI:    "spotify:user:attributes:update",
		Decode: DecodePB[*parental_controls.UserAttributesUpdate],
	}
	TopicUserAttributesMutated = Topic[*parental_controls.UserAttributesUpdate]{
		URI:    "spotify:user:attributes:mutated",
		Decode: DecodePB[*parental_controls.UserAttributesUpdate],
	}
	TopicDeviceSettingsChanged = Topic[*pubsub.DeviceSettingsFieldsChangedEvent]{
		URI:    "hm://playback/v1/devicesettings",
		Decode: DecodePB[*pubsub.DeviceSettingsFieldsChangedEvent],
	}
	TopicDeviceBroadcastStatus = Topic[*types.DeviceBroadcastStatus]{
		URI:    "social-connect/v2/broadcast_status_updates",
		Decode: decodeDeviceBroadcastStatus,
	}
)
