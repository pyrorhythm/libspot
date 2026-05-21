package connect

import (
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

// parseState decodes a connect-state JSON body into the public State snapshot.
func parseState(body []byte) (State, error) {
	val, err := fastjson.ParseBytes(body)
	if err != nil {
		return State{}, errors.Wrap(err, "connect: parse state")
	}

	state := State{
		ActiveDeviceID: string(val.Get("active_device_id").GetStringBytes()),
	}

	if devices := val.Get("devices"); devices != nil {
		state.Devices = parseDevices(devices, state.ActiveDeviceID)
	}
	if state.ActiveDeviceID == "" {
		state.ActiveDeviceID = detectActiveDeviceID(val.Get("devices"))
	}
	if player := val.Get("player_state"); player != nil {
		state.Player = parsePlayerState(player)
		state.OriginDeviceID = playOriginDeviceID(player)
	}

	return state, nil
}

func parseDevices(devices *fastjson.Value, activeID string) []Device {
	if devices == nil || devices.Type() != fastjson.TypeObject {
		return nil
	}
	out := make([]Device, 0)
	devices.GetObject().Visit(func(key []byte, v *fastjson.Value) {
		id := string(key)
		out = append(out, parseDevice(id, v, id == activeID))
	})
	return out
}

func parseDevice(id string, v *fastjson.Value, active bool) Device {
	device := Device{ID: id, Active: active}
	if v == nil || v.Type() != fastjson.TypeObject {
		return device
	}
	device.Name = string(v.Get("name").GetStringBytes())
	if device.Name == "" {
		device.Name = string(v.Get("device_name").GetStringBytes())
	}
	device.Type = string(v.Get("device_type").GetStringBytes())
	device.Volume = normalizeConnectVolume(
		v.Get("volume").GetInt(),
		v.Get("volume_percent").GetInt(),
	)
	return device
}

func detectActiveDeviceID(devices *fastjson.Value) string {
	if devices == nil || devices.Type() != fastjson.TypeObject {
		return ""
	}
	var found string
	devices.GetObject().Visit(func(key []byte, v *fastjson.Value) {
		if found != "" || v == nil {
			return
		}
		if v.Get("is_active").GetBool() ||
			v.Get("is_currently_playing").GetBool() ||
			v.Get("is_active_device").GetBool() {
			found = string(key)
		}
	})
	return found
}

func playOriginDeviceID(player *fastjson.Value) string {
	if player == nil {
		return ""
	}
	return string(player.Get("play_origin", "device_identifier").GetStringBytes())
}

func parsePlayerState(player *fastjson.Value) *PlayerState {
	if player == nil || player.Type() != fastjson.TypeObject {
		return nil
	}
	ps := &PlayerState{
		Shuffle: player.Get("shuffle").GetBool(),
		Repeat:  string(player.Get("repeat_mode").GetStringBytes()),
	}
	if ps.Repeat == "" {
		ps.Repeat = string(player.Get("repeat").GetStringBytes())
	}
	if player.Get("is_paused").Exists() {
		ps.IsPlaying = !player.Get("is_paused").GetBool()
	} else {
		ps.IsPlaying = player.Get("is_playing").GetBool()
	}
	ps.ProgressMS = player.Get("position_as_of_timestamp").GetInt()
	if ps.ProgressMS == 0 {
		ps.ProgressMS = player.Get("position_ms").GetInt()
	}
	ps.NowPlaying = parseNowPlaying(player)
	ps.UpNext = parseUpNext(player.Get("next_tracks"))
	return ps
}

func parseNowPlaying(player *fastjson.Value) *MediaOneof {
	for _, key := range []string{"track", "item", "current_track"} {
		if raw := player.Get(key); raw != nil {
			if o, ok := parseMediaOneof(raw); ok {
				return o
			}
		}
	}
	return nil
}

func parseUpNext(next *fastjson.Value) []MediaOneof {
	if next == nil {
		return nil
	}
	arr, err := next.Array()
	if err != nil {
		return nil
	}
	out := make([]MediaOneof, 0, len(arr))
	for _, entry := range arr {
		if o, ok := parseMediaOneof(entry); ok {
			out = append(out, *o)
		}
	}
	return out
}

func parseMediaOneof(v *fastjson.Value) (*MediaOneof, bool) {
	if v == nil || v.Type() != fastjson.TypeObject {
		return nil, false
	}
	uri := string(v.Get("uri").GetStringBytes())
	if uri == "" {
		uri = string(v.Get("metadata", "uri").GetStringBytes())
	}
	if uri == "" {
		return nil, false
	}
	name := string(v.Get("name").GetStringBytes())
	album := string(v.Get("metadata", "album_title").GetStringBytes())
	var artists []string
	if artist := string(v.Get("metadata", "artist_name").GetStringBytes()); artist != "" {
		artists = []string{artist}
	}
	if name == "" {
		name = string(v.Get("metadata", "title").GetStringBytes())
	}
	if typeFromURI(uri) == string(mediaEpisode) {
		return PartialEpisode(idFromURI(uri), uri, name), true
	}
	return PartialTrack(idFromURI(uri), uri, name, album, artists), true
}

func mapPlayback(state State) Playback {
	pb := Playback{}
	if state.Player != nil {
		pb.IsPlaying = state.Player.IsPlaying
		pb.ProgressMS = state.Player.ProgressMS
		pb.Shuffle = state.Player.Shuffle
		pb.Repeat = state.Player.Repeat
		if state.Player.NowPlaying != nil {
			item := state.Player.NowPlaying.Item()
			pb.NowPlaying = &item
		}
	}
	for _, d := range state.Devices {
		if d.Active {
			pb.Device = d
			break
		}
	}
	return pb
}

func mapQueue(state State) Queue {
	q := Queue{}
	if state.Player == nil {
		return q
	}
	q.CurrentlyPlaying = state.Player.NowPlaying
	q.UpNext = state.Player.UpNext
	return q
}

// normalizeConnectVolume converts Spotify's 0..65535 internal scale to 0..100,
// leaving values that already look like percentages untouched.
func normalizeConnectVolume(values ...int) int {
	for _, volume := range values {
		if volume == 0 {
			continue
		}
		if volume <= 100 {
			return clampVolume(volume)
		}
		return clampVolume(int(float64(volume)*100/65535 + 0.5))
	}
	return 0
}
