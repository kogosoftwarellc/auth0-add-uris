package types

type RequestBodyCreateToken struct {
	GrantType string `json:"grant_type"`
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience string `json:"audience"`
}
