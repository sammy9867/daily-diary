package domain

// TokenDetail model struct
type TokenDetail struct {
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	AccessUUID        string `json:"-"`
	RefreshUUID       string `json:"-"`
	AcctokenExpiresAt int64  `json:"-"`
	ReftokenExpiresAt int64  `json:"-"`
}
