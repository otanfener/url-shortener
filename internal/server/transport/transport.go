package transport

type ShortenURLRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenURLResponse struct {
	ShortCode string `json:"short_code"`
}
