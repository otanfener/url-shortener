package dto

type ShortenRequest struct {
	LongURL string
}

type ShortenResponse struct {
	ShortCode string
}
