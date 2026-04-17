package pfdomain

type (
	WatchFeedEntrypoint struct {
		EntrypointURI  string         `json:"entrypointUri"`
		ThumbnailImage ThumbnailImage `json:"thumbnailImage"`
		Video          Video          `json:"video"`
	}

	Video struct {
		VideoType VideoType `json:"videoType"`
		FileID    string    `json:"fileId"`

		EndTime   int `json:"endTime"`
		StartTime int `json:"startTime"`
	}
)
