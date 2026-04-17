package pfdomain

type Color struct {
	Hex        string `json:"hex"`
	IsFallback bool   `json:"isFallback"`
}

type ExtractedColor struct {
	ColorDark  *Color `json:"colorDark,omitempty"`
	ColorLight *Color `json:"colorLight,omitempty"`
	ColorRaw   *Color `json:"colorRaw,omitempty"`
}

type RGBA struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
	Alpha int `json:"alpha"`
}

type ExtractedColorExtended struct {
	BackgroundBase       RGBA `json:"backgroundBase"`
	BackgroundTintedBase RGBA `json:"backgroundTintedBase"`
	TextBase             RGBA `json:"textBase"`
	TextBrightAccent     RGBA `json:"textBrightAccent"`
	TextSubdued          RGBA `json:"textSubdued"`
}

type ExtractedColorSet struct {
	EncoreBaseSetTextColor *RGBA                   `json:"encoreBaseSetTextColor"`
	HighContrast           *ExtractedColorExtended `json:"highContrast"`
	HigherContrast         *ExtractedColorExtended `json:"higherContrast"`
	MinContrast            *ExtractedColorExtended `json:"minContrast"`
}

type DynamicColorSet struct {
	Status  string             `json:"status"` // "OK" or ....
	BestFit string             `json:"bestFit"`
	Dark    *ExtractedColorSet `json:"dark"`
	Light   *ExtractedColorSet `json:"light"`
}
