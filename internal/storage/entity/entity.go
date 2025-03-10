package entity

type URLMapping struct {
	ShortCode string `dynamodbav:"short_code"`
	LongURL   string `dynamodbav:"long_url"`
}
