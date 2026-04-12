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
