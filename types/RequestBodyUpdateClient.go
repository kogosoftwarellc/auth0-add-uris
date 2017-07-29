package types

type RequestBodyUpdateClient struct {
	Callbacks []string `json:"callbacks"`
	AllowedLogoutUrls []string `json:"allowed_logout_urls"`
}
