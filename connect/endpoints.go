package connect

type Endpoint string

const (
	EndpointSetOptions          Endpoint = "set_options"
	EndpointResume              Endpoint = "resume"
	EndpointPlay                Endpoint = "play"
	EndpointPause               Endpoint = "pause"
	EndpointSkipNext            Endpoint = "skip_next"
	EndpointSkipPrev            Endpoint = "skip_prev"
	EndpointSeekTo              Endpoint = "seek_to"
	EndpointSetShufflingContext Endpoint = "set_shuffling_context"
	EndpointAddToQueue          Endpoint = "add_to_queue"
)

// RepeatMode is the repeat mode sent via set_options.
type RepeatMode string

const (
	RepeatOff     RepeatMode = ""
	RepeatTrack   RepeatMode = "track"
	RepeatContext RepeatMode = "context"
)
