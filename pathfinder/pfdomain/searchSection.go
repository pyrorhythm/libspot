package pfdomain

type SectionTitle struct {
	TransformedLabel string `json:"transformedLabel"`
}

type SearchSection struct {
	AttributionUri string        `json:"attributionUri"`
	SectionType    string        `json:"sectionType"`
	Title          *SectionTitle `json:"title"`
	Items          []*Oneof      `json:"items"`
}
