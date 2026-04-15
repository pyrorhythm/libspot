package domaintypes

type AssociationsV3 struct {
	AudioAssociations TotalCount `json:"audioAssociations"`
	VideoAssociations TotalCount `json:"videoAssociations"`
}
type TotalCount struct {
	TotalCount int `json:"totalCount"`
}
