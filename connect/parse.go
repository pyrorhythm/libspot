package connect

import (
	"time"

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
	if player := extractConnectPlayerJSON(val); player != nil {
		state.Player = parseConnectPlayer(player)
		state.OriginDeviceID = playOriginDeviceID(player)
	}

	return state, nil
}

// extractConnectPlayerJSON locates the connect player snapshot in a connect-state body.
// Spotify may return it at the root (track, context_uri, …) or nested under player_state / cluster.
func extractConnectPlayerJSON(val *fastjson.Value) *fastjson.Value {
	if val == nil {
		return nil
	}
	if ps := val.Get("player_state"); ps != nil && ps.Type() == fastjson.TypeObject {
		return ps
	}
	if cluster := val.Get("cluster"); cluster != nil {
		if ps := cluster.Get("player_state"); ps != nil && ps.Type() == fastjson.TypeObject {
			return ps
		}
	}
	if looksLikeConnectPlayerJSON(val) {
		return val
	}
	return nil
}

func looksLikeConnectPlayerJSON(v *fastjson.Value) bool {
	if v == nil || v.Type() != fastjson.TypeObject {
		return false
	}
	if track := v.Get("track"); track != nil && track.Type() == fastjson.TypeObject {
		return true
	}
	return v.Get("context_uri").Exists() &&
		(v.Get("is_playing").Exists() || v.Get("is_paused").Exists() || v.Get("position_as_of_timestamp").Exists())
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
		jsonInt(v.Get("volume")),
		jsonInt(v.Get("volume_percent")),
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

func parseConnectPlayer(player *fastjson.Value) *PlayerState {
	if player == nil || player.Type() != fastjson.TypeObject {
		return nil
	}

	ps := &PlayerState{
		IsPlaying:  parseIsPlaying(player),
		Shuffle:    parseShuffle(player),
		Repeat:     parseRepeat(player),
		ProgressMS: extrapolateProgress(player),
		DurationMS: jsonInt(player.Get("duration")),
		ContextURI: string(player.Get("context_uri").GetStringBytes()),
		NowPlaying: parseNowPlaying(player),
		UpNext:     parseUpNext(player.Get("next_tracks")),
	}
	return ps
}

// parseIsPlaying treats is_paused as authoritative when present (Spotify may send both flags).
func parseIsPlaying(player *fastjson.Value) bool {
	if player.Get("is_paused").Exists() {
		return !player.Get("is_paused").GetBool()
	}
	return player.Get("is_playing").GetBool()
}

func parseShuffle(player *fastjson.Value) bool {
	if player.Get("shuffle").Exists() {
		return player.Get("shuffle").GetBool()
	}
	return player.Get("options", "shuffling_context").GetBool()
}

func parseRepeat(player *fastjson.Value) string {
	if r := string(player.Get("repeat_mode").GetStringBytes()); r != "" {
		return r
	}
	if r := string(player.Get("repeat").GetStringBytes()); r != "" {
		return r
	}
	opts := player.Get("options")
	if opts == nil {
		return ""
	}
	if opts.Get("repeating_track").GetBool() {
		return "track"
	}
	if opts.Get("repeating_context").GetBool() {
		return "context"
	}
	return "off"
}

func extrapolateProgress(player *fastjson.Value) int {
	pos := jsonInt(player.Get("position_as_of_timestamp"))
	if pos == 0 {
		pos = jsonInt(player.Get("position_ms"))
	}
	if !parseIsPlaying(player) {
		return pos
	}
	ts := jsonInt64(player.Get("timestamp"))
	if ts <= 0 {
		return pos
	}
	now := time.Now().UnixMilli()
	if now <= ts {
		return pos
	}
	pos += int(now - ts)
	if dur := jsonInt(player.Get("duration")); dur > 0 && pos > dur {
		return dur
	}
	return pos
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
		if isSkippedQueueEntry(entry) {
			continue
		}
		if o, ok := parseMediaOneof(entry); ok {
			out = append(out, *o)
		}
	}
	return out
}

func isSkippedQueueEntry(v *fastjson.Value) bool {
	if v == nil {
		return true
	}
	uri := string(v.Get("uri").GetStringBytes())
	if uri == "" {
		return true
	}
	kind := typeFromURI(uri)
	if kind != string(mediaTrack) && kind != string(mediaEpisode) {
		return true
	}
	if string(v.Get("metadata", "hidden").GetStringBytes()) == "true" {
		return true
	}
	return false
}

func parseMediaOneof(v *fastjson.Value) (*MediaOneof, bool) {
	if v == nil || v.Type() != fastjson.TypeObject {
		return nil, false
	}
	uri := string(v.Get("uri").GetStringBytes())
	if uri == "" {
		uri = string(v.Get("metadata", "uri").GetStringBytes())
	}
	if uri == "" || isSkippedQueueEntry(v) {
		return nil, false
	}
	meta := v.Get("metadata")
	name := string(v.Get("name").GetStringBytes())
	if name == "" {
		name = string(meta.Get("title").GetStringBytes())
	}
	album := string(meta.Get("album_title").GetStringBytes())
	artists := parseArtists(meta)
	if typeFromURI(uri) == string(mediaEpisode) {
		return PartialEpisode(idFromURI(uri), uri, name), true
	}
	return PartialTrack(idFromURI(uri), uri, name, album, artists), true
}

func parseArtists(meta *fastjson.Value) []string {
	if meta == nil {
		return nil
	}
	if artist := string(meta.Get("artist_name").GetStringBytes()); artist != "" {
		return []string{artist}
	}
	return nil
}

func mapPlayback(state State) Playback {
	pb := Playback{}
	if state.Player != nil {
		pb.IsPlaying = state.Player.IsPlaying
		pb.ProgressMS = state.Player.ProgressMS
		pb.DurationMS = state.Player.DurationMS
		pb.Shuffle = state.Player.Shuffle
		pb.Repeat = state.Player.Repeat
		pb.ContextURI = state.Player.ContextURI
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
