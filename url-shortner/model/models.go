package model

type User struct {
	ID       int64  `bson:"_id"`
	LongUrl  string `bson:"long_url"`
	ShortUrl string `bson:"short_url"`
	Clicks   int    `bson:"clicks"`
}
