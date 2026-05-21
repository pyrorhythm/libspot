package connect

// LoggingParams mirrors connect-state command logging_params.
type LoggingParams struct {
	CommandID string `json:"command_id"`
}

func newLoggingParams() LoggingParams {
	return LoggingParams{CommandID: randomID()}
}

// SimpleCommandPayload is the body for endpoint-only player commands.
type SimpleCommandPayload struct {
	Command SimplePlayerCommand `json:"command"`
}

type SimplePlayerCommand struct {
	Endpoint      string        `json:"endpoint"`
	LoggingParams LoggingParams `json:"logging_params"`
}

func newSimpleCommand(endpoint string) SimpleCommandPayload {
	return SimpleCommandPayload{
		Command: SimplePlayerCommand{
			Endpoint:      endpoint,
			LoggingParams: newLoggingParams(),
		},
	}
}

// SeekCommandPayload is the body for seek_to.
type SeekCommandPayload struct {
	Command SeekPlayerCommand `json:"command"`
}

type SeekPlayerCommand struct {
	Endpoint      string        `json:"endpoint"`
	Value         int           `json:"value"`
	LoggingParams LoggingParams `json:"logging_params"`
}

func newSeekCommand(positionMS int) SeekCommandPayload {
	return SeekCommandPayload{
		Command: SeekPlayerCommand{
			Endpoint:      "seek_to",
			Value:         positionMS,
			LoggingParams: newLoggingParams(),
		},
	}
}

// ShuffleCommandPayload is the body for set_shuffling_context.
type ShuffleCommandPayload struct {
	Command ShufflePlayerCommand `json:"command"`
}

type ShufflePlayerCommand struct {
	Endpoint      string        `json:"endpoint"`
	Value         bool          `json:"value"`
	LoggingParams LoggingParams `json:"logging_params"`
}

func newShuffleCommand(enabled bool) ShuffleCommandPayload {
	return ShuffleCommandPayload{
		Command: ShufflePlayerCommand{
			Endpoint:      "set_shuffling_context",
			Value:         enabled,
			LoggingParams: newLoggingParams(),
		},
	}
}

// SetOptionsCommandPayload is the body for set_options (repeat mode).
type SetOptionsCommandPayload struct {
	Command SetOptionsPlayerCommand `json:"command"`
}

type SetOptionsPlayerCommand struct {
	Endpoint          string        `json:"endpoint"`
	RepeatingTrack    bool          `json:"repeating_track"`
	RepeatingContext  bool          `json:"repeating_context"`
	LoggingParams     LoggingParams `json:"logging_params"`
}

func newSetOptionsCommand(repeatingTrack, repeatingContext bool) SetOptionsCommandPayload {
	return SetOptionsCommandPayload{
		Command: SetOptionsPlayerCommand{
			Endpoint:         "set_options",
			RepeatingTrack:   repeatingTrack,
			RepeatingContext: repeatingContext,
			LoggingParams:    newLoggingParams(),
		},
	}
}

// AddToQueueCommandPayload is the body for add_to_queue.
type AddToQueueCommandPayload struct {
	Command AddToQueuePlayerCommand `json:"command"`
}

type QueueTrackRef struct {
	URI string `json:"uri"`
}

type AddToQueuePlayerCommand struct {
	Endpoint      string        `json:"endpoint"`
	Track         QueueTrackRef `json:"track"`
	LoggingParams LoggingParams `json:"logging_params"`
}

func newAddToQueueCommand(uri string) AddToQueueCommandPayload {
	return AddToQueueCommandPayload{
		Command: AddToQueuePlayerCommand{
			Endpoint:      "add_to_queue",
			Track:         QueueTrackRef{URI: uri},
			LoggingParams: newLoggingParams(),
		},
	}
}

// PlayCommandPayload is the body for play with a context or track uri.
type PlayCommandPayload struct {
	Command PlayPlayerCommand `json:"command"`
}

type PlayContextRef struct {
	URI string `json:"uri"`
	URL string `json:"url"`
}

type PlaySkipTo struct {
	TrackURI string `json:"track_uri"`
}

type PlayCommandOptions struct {
	SkipTo PlaySkipTo `json:"skip_to"`
}

type PlayPlayerCommand struct {
	Endpoint      string              `json:"endpoint"`
	LoggingParams LoggingParams       `json:"logging_params"`
	Context       PlayContextRef      `json:"context"`
	Options       *PlayCommandOptions `json:"options,omitempty"`
}

func newPlayCommand(uri string) PlayCommandPayload {
	cmd := PlayPlayerCommand{
		Endpoint:      "play",
		LoggingParams: newLoggingParams(),
		Context: PlayContextRef{
			URI: uri,
			URL: "context://" + uri,
		},
	}
	if !IsContextURI(uri) {
		cmd.Options = &PlayCommandOptions{
			SkipTo: PlaySkipTo{TrackURI: uri},
		}
	}
	return PlayCommandPayload{Command: cmd}
}

// VolumePayload is the body for PUT .../connect/volume/...
type VolumePayload struct {
	Volume int `json:"volume"`
}

// TransferPayload is the body for POST .../connect/transfer/...
type TransferPayload struct {
	TransferOptions TransferOptions `json:"transfer_options"`
	CommandID       string          `json:"command_id"`
}

type TransferOptions struct {
	RestorePaused string `json:"restore_paused"`
}

func newTransferPayload() TransferPayload {
	return TransferPayload{
		TransferOptions: TransferOptions{RestorePaused: "resume"},
		CommandID:       randomID(),
	}
}

// StateRequestPayload is the body for PUT .../devices/hobs_* (connect-state poll).
type StateRequestPayload struct {
	MemberType string              `json:"member_type"`
	Device     StateRequestDevice  `json:"device"`
}

type StateRequestDevice struct {
	DeviceInfo StateRequestDeviceInfo `json:"device_info"`
}

type StateRequestDeviceInfo struct {
	Capabilities StateRequestCapabilities `json:"capabilities"`
}

type StateRequestCapabilities struct {
	CanBePlayer          bool `json:"can_be_player"`
	Hidden               bool `json:"hidden"`
	NeedsFullPlayerState bool `json:"needs_full_player_state"`
}

func newStateRequestPayload() StateRequestPayload {
	return StateRequestPayload{
		MemberType: "CONNECT_STATE",
		Device: StateRequestDevice{
			DeviceInfo: StateRequestDeviceInfo{
				Capabilities: StateRequestCapabilities{
					CanBePlayer:          false,
					Hidden:               true,
					NeedsFullPlayerState: true,
				},
			},
		},
	}
}

// RegisterDevicePayload is the body for POST .../track-playback/v1/devices.
type RegisterDevicePayload struct {
	Device                  RegisterDeviceBody `json:"device"`
	OutroEndcontentSnooping bool                 `json:"outro_endcontent_snooping"`
	ConnectionID            string               `json:"connection_id"`
	ClientVersion           string               `json:"client_version"`
	Volume                  int                  `json:"volume"`
}

type RegisterDeviceBody struct {
	DeviceID           string                    `json:"device_id"`
	DeviceType         string                    `json:"device_type"`
	Brand              string                    `json:"brand"`
	Model              string                    `json:"model"`
	Name               string                    `json:"name"`
	IsGroup            bool                      `json:"is_group"`
	Metadata           struct{}                  `json:"metadata"`
	PlatformIdentifier string                    `json:"platform_identifier"`
	Capabilities       RegisterDeviceCapabilities `json:"capabilities"`
}

type RegisterDeviceCapabilities struct {
	ChangeVolume           bool     `json:"change_volume"`
	SupportsFileMediaType  bool     `json:"supports_file_media_type"`
	EnablePlayToken        bool     `json:"enable_play_token"`
	PlayTokenLostBehavior  string   `json:"play_token_lost_behavior"`
	DisableConnect         bool     `json:"disable_connect"`
	AudioPodcasts          bool     `json:"audio_podcasts"`
	VideoPlayback          bool     `json:"video_playback"`
	ManifestFormats        []string `json:"manifest_formats"`
}

func newRegisterDevicePayload(deviceID, connectionID, platformID string) RegisterDevicePayload {
	return RegisterDevicePayload{
		Device: RegisterDeviceBody{
			DeviceID:           deviceID,
			DeviceType:         "computer",
			Brand:              "spotify",
			Model:              connectDeviceModel,
			Name:               connectDeviceName,
			IsGroup:            false,
			PlatformIdentifier: platformID,
			Capabilities: RegisterDeviceCapabilities{
				ChangeVolume:          true,
				SupportsFileMediaType: true,
				EnablePlayToken:       true,
				PlayTokenLostBehavior: "pause",
				DisableConnect:        false,
				AudioPodcasts:         true,
				VideoPlayback:         true,
				ManifestFormats: []string{
					"file_ids_mp3",
					"file_urls_mp3",
					"file_ids_mp4",
					"manifest_ids_video",
				},
			},
		},
		OutroEndcontentSnooping: false,
		ConnectionID:            connectionID,
		ClientVersion:           clientVersion,
		Volume:                  65535,
	}
}
