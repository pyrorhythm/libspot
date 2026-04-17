package pfdomain

type Genre struct {
	URI   string `json:"uri"`
	Name  string `json:"name"`
	Image *Image `json:"image"`
}
