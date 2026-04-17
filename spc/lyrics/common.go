package lyrics

const Path = "color-lyrics/v2/track/{trackB62Id}/image/spotify%3Aimage%3A{imageGid}"

type LyricLine struct {
	StartTimeMs         string   `json:"startTimeMs"`
	Words               string   `json:"words"`
	Syllables           []string `json:"syllables"`
	EndTimeMs           string   `json:"endTimeMs"`
	TransliteratedWords string   `json:"transliteratedWords"`
}

type Lyrics struct {
	SyncType            string      `json:"syncType"`
	Lines               []LyricLine `json:"lines"`
	Provider            string      `json:"provider"`
	ProviderLyricsId    string      `json:"providerLyricsId"`
	ProviderDisplayName string      `json:"providerDisplayName"`
	SyncLyricsUri       string      `json:"syncLyricsUri"`
	IsDenseTypeface     bool        `json:"isDenseTypeface"`
	Alternatives        []any       `json:"alternatives"`
	Language            string      `json:"language"`
	IsRtlLanguage       bool        `json:"isRtlLanguage"`
	CapStatus           string      `json:"capStatus"`
	PreviewLines        []LyricLine `json:"previewLines"`
}

type Colors struct {
	Background    int `json:"background"`
	Text          int `json:"text"`
	HighlightText int `json:"highlightText"`
}

type Response struct {
	Lyrics          Lyrics `json:"lyrics"`
	Colors          Colors `json:"colors"`
	HasVocalRemoval bool   `json:"hasVocalRemoval"`
}
