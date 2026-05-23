package connect

// State is the decoded connect-state snapshot returned by the connect-state API.
type State struct {
	ActiveDeviceID string
	OriginDeviceID string
	Devices        []Device
	Player         *PlayerState
}

// Device is a Spotify Connect endpoint visible in the session.
type Device struct {
	ID     string
	Name   string
	Type   string
	Active bool
	Volume int // 0..100
}

// PlayerState is the live playback snapshot from connect-state.
type PlayerState struct {
	IsPlaying  bool
	ProgressMS int
	DurationMS int
	Shuffle    bool
	Repeat     string
	NowPlaying *MediaOneof
	UpNext     []MediaOneof
}

// Playback mirrors the fields callers typically need from PlayerState plus the
// active device.
type Playback struct {
	IsPlaying  bool
	ProgressMS int
	DurationMS int
	Shuffle    bool
	Repeat     string
	NowPlaying *MediaItem
	Device     Device
}

// Queue is the currently playing item plus upcoming tracks.
type Queue struct {
	CurrentlyPlaying *MediaOneof
	UpNext           []MediaOneof
}

// MediaItem is a flattened view of a track or episode for simple consumers.
type MediaItem struct {
	ID      string
	URI     string
	Name    string
	Type    string
	Album   string
	Artists []string
}
