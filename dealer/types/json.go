package types

type DeviceBroadcastStatus struct {
	Timestamp        string `json:"timestamp"`
	BroadcastStatus  string `json:"broadcast_status"`
	DeviceId         string `json:"device_id"`
	OutputDeviceInfo struct {
		OutputDeviceType string `json:"output_device_type"`
		DeviceName       string `json:"device_name"`
	} `json:"output_device_info"`
	LinkToken struct {
		Token string `json:"token"`
	} `json:"link_token"`
	DeviceType string `json:"device_type"`
}

type SessionUpdate struct {
	Session struct {
		Timestamp        string `json:"timestamp"`
		SessionId        string `json:"session_id"`
		JoinSessionToken string `json:"join_session_token"`
		JoinSessionUrl   string `json:"join_session_url"`
		SessionOwnerId   string `json:"session_owner_id"`
		SessionMembers   []struct {
			JoinedTimestamp string `json:"joined_timestamp"`
			Id              string `json:"id"`
			Username        string `json:"username"`
			DisplayName     string `json:"display_name"`
			ImageUrl        string `json:"image_url"`
			LargeImageUrl   string `json:"large_image_url"`
			IsListening     bool   `json:"is_listening"`
			IsControlling   bool   `json:"is_controlling"`
			PlaybackControl string `json:"playbackControl"`
			IsCurrentUser   bool   `json:"is_current_user"`
		} `json:"session_members"`
		JoinSessionUri     string `json:"join_session_uri"`
		IsSessionOwner     bool   `json:"is_session_owner"`
		IsListening        bool   `json:"is_listening"`
		IsControlling      bool   `json:"is_controlling"`
		InitialSessionType string `json:"initialSessionType"`
		HostActiveDeviceId string `json:"hostActiveDeviceId"`
		MaxMemberCount     int    `json:"maxMemberCount"`
		QueueOnlyMode      bool   `json:"queue_only_mode"`
		WifiBroadcast      bool   `json:"wifi_broadcast"`
		HostDeviceInfo     struct {
			DeviceId         string `json:"device_id"`
			OutputDeviceInfo struct {
				OutputDeviceType string `json:"output_device_type"`
				DeviceName       string `json:"device_name"`
			} `json:"output_device_info"`
			DeviceName string `json:"device_name"`
			DeviceType string `json:"device_type"`
			IsGroup    bool   `json:"is_group"`
		} `json:"host_device_info"`
		IsPaused bool `json:"is_paused"`
	} `json:"session"`
	Reason               string `json:"reason"`
	UpdateSessionMembers []struct {
		JoinedTimestamp string `json:"joined_timestamp"`
		Id              string `json:"id"`
		Username        string `json:"username"`
		DisplayName     string `json:"display_name"`
		ImageUrl        string `json:"image_url"`
		LargeImageUrl   string `json:"large_image_url"`
		IsListening     bool   `json:"is_listening"`
		IsControlling   bool   `json:"is_controlling"`
		PlaybackControl string `json:"playbackControl"`
	} `json:"update_session_members"`
}
