package pfdomain

type ColorDark struct {
	Hex        string `json:"hex"`
	IsFallback bool   `json:"isFallback"`
}

type RGBA struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
	Alpha int `json:"alpha"`
}

type ExtractedColor struct {
	BackgroundBase       RGBA `json:"backgroundBase"`
	BackgroundTintedBase RGBA `json:"backgroundTintedBase"`
	TextBase             RGBA `json:"textBase"`
	TextBrightAccent     RGBA `json:"textBrightAccent"`
	TextSubdued          RGBA `json:"textSubdued"`
}

type ExtractedColorSet struct {
	EncoreBaseSetTextColor *RGBA           `json:"encoreBaseSetTextColor"`
	HighContrast           *ExtractedColor `json:"highContrast"`
	HigherContrast         *ExtractedColor `json:"higherContrast"`
	MinContrast            *ExtractedColor `json:"minContrast"`
}

type VisualIdentity struct {
	SixteenByNineCoverImage *VisualIdentityImage `json:"sixteenByNineCoverImage"`
	SquareCoverImage        *VisualIdentityImage `json:"squareCoverImage"`
}
