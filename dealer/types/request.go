package types

import (
	"github.com/pyrorhythm/libspot/gen/spotify/connectstate"
	"github.com/pyrorhythm/libspot/gen/spotify/player"
)

type Request struct {
	resp chan<- bool

	Key     string
	Ident   string
	Uri     string
	Method  string
	Headers map[string]string
	Payload RequestPayload
}

// BindResponder wires the reply channel. Called by the dealer connection
// after dispatching the request to a handler; user code should not invoke it.
func (r *Request) BindResponder(ch chan<- bool) {
	r.resp = ch
}

func (r *Request) Respond(ok bool) {
	if r.resp == nil {
		return
	}
	select {
	case r.resp <- ok:
	default:
	}
}

type RequestPayload struct {
	MessageId      uint32  `json:"message_id"`
	SentByDeviceId string  `json:"sent_by_device_id"`
	Command        Command `json:"command"`
}

type Command struct {
	Endpoint             string                   `json:"endpoint"`
	SessionId            string                   `json:"session_id"`
	Data                 []byte                   `json:"data"`
	Value                interface{}              `json:"value"`
	Position             int64                    `json:"position"`
	Relative             string                   `json:"relative"`
	Context              *player.Context          `json:"context"`
	PlayOrigin           *connectstate.PlayOrigin `json:"play_origin"`
	Track                *player.ContextTrack     `json:"track"`
	PrevTracks           []*player.ContextTrack   `json:"prev_tracks"`
	NextTracks           []*player.ContextTrack   `json:"next_tracks"`
	RepeatingTrack       *bool                    `json:"repeating_track"`
	RepeatingContext     *bool                    `json:"repeating_context"`
	ShufflingContext     *bool                    `json:"shuffling_context"`
	LoggingParams        LoggingParams            `json:"logging_params"`
	Options              Options                  `json:"options"`
	PlayOptions          PlayOptions              `json:"play_options"`
	FromDeviceIdentifier string                   `json:"from_device_identifier"`
}

type LoggingParams struct {
	CommandInitiatedTime int64    `json:"command_initiated_time"`
	PageInstanceIds      []string `json:"page_instance_ids"`
	InteractionIds       []string `json:"interaction_ids"`
	DeviceIdentifier     string   `json:"device_identifier"`
}

type Options struct {
	RestorePaused         string                               `json:"restore_paused"`
	RestorePosition       string                               `json:"restore_position"`
	RestoreTrack          string                               `json:"restore_track"`
	AlwaysPlaySomething   bool                                 `json:"always_play_something"`
	AllowSeeking          bool                                 `json:"allow_seeking"`
	SkipTo                SkipTo                               `json:"skip_to"`
	InitiallyPaused       bool                                 `json:"initially_paused"`
	SystemInitiated       bool                                 `json:"system_initiated"`
	PlayerOptionsOverride *player.ContextPlayerOptionOverrides `json:"player_options_override"`
	Suppressions          *connectstate.Suppressions           `json:"suppressions"`
	PrefetchLevel         string                               `json:"prefetch_level"`
	AudioStream           string                               `json:"audio_stream"`
	SessionId             string                               `json:"session_id"`
	License               string                               `json:"license"`
}

type SkipTo struct {
	TrackUid   string `json:"track_uid"`
	TrackUri   string `json:"track_uri"`
	TrackIndex int    `json:"track_index"`
}

type PlayOptions struct {
	OverrideRestrictions bool   `json:"override_restrictions"`
	OnlyForLocalDevice   bool   `json:"only_for_local_device"`
	SystemInitiated      bool   `json:"system_initiated"`
	Reason               string `json:"reason"`
	Operation            string `json:"operation"`
	Trigger              string `json:"trigger"`
}

type Reply struct {
	Type    string `json:"type"`
	Key     string `json:"key"`
	Payload struct {
		Success bool `json:"success"`
	} `json:"payload"`
}
