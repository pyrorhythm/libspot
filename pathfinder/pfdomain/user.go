package pfdomain

type (
	User struct {
		Uri      string   `json:"uri"`
		Name     string   `json:"name"`
		Username string   `json:"username"`
		Avatar   ImageRaw `json:"avatar"`
	}

	IsFollowingUser struct {
		Uri       string `json:"uri"`
		Following bool   `json:"following"`
	}
)
