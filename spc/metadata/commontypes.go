package metadata

type (
	Copyright struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}

	OriginalAudio struct {
		Uuid   string `json:"uuid"`
		Format string `json:"format"`
	}

	Uuid struct {
		Uuid string `json:"uuid"`
	}

	ExternalId struct {
		Type string `json:"type"`
		Id   string `json:"id"`
	}

	ArtistShort struct {
		Gid  string `json:"gid"`
		Name string `json:"name"`
	}

	Image struct {
		FileId string `json:"file_id"`
		Size   string `json:"size"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}

	CoverGroup struct {
		Image []Image `json:"image"`
	}

	Date struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	}

	Unix struct {
		Seconds int `json:"seconds"`
		Nanos   int `json:"nanos"`
	}

	ImplementationDetails struct {
		CatalogInsertionDate Unix `json:"catalog_insertion_date"`
	}

	Gid struct {
		Gid string `json:"gid"`
	}
)
