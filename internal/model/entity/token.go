package entity

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewTokens(access, refresh string) *Tokens {
	return &Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}
