package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type NewsItem struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Path    string             `bson:"path,omitempty" json:"path,omitempty"`
	Time    int64              `bson:"time,omitempty" json:"time,omitempty"`
	Status  string             `bson:"status,omitempty" json:"status,omitempty"`
	Content []NewsItemContent  `bson:"content,omitempty" json:"content,omitempty"`
}

type NewsItemContent struct {
	Language    string       `bson:"language,omitempty" json:"language,omitempty"`
	Title       string       `bson:"title,omitempty" json:"title,omitempty"`
	Image       Image        `bson:"image,omitempty" json:"image,omitempty"`
	MarkdownURL string       `bson:"markdown,omitempty" json:"markdown,omitempty"`
	Links       []Link       `bson:"links,omitempty" json:"links,omitempty"`
	Overview    OverviewInfo `bson:"overview,omitempty" json:"overview,omitempty"`
}

type OverviewInfo struct {
	Title       string `bson:"title,omitempty" json:"title,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Image       Image  `bson:"image,omitempty" json:"image,omitempty"`
}

type Image struct {
	URL       string `bson:"url,omitempty" json:"url,omitempty"`
	Copyright string `bson:"copyright,omitempty" json:"copyright,omitempty"`
}

type Link struct {
	Title         string `bson:"title,omitempty" json:"title,omitempty"`
	Description   string `bson:"description,omitempty" json:"description,omitempty"`
	Image         Image  `bson:"image,omitempty" json:"image,omitempty"`
	ImagePosition string `bson:"imagePosition,omitempty" json:"imagePosition,omitempty"` //e.g. top left background
	LinkText      string `bson:"linkText,omitempty" json:"linkText,omitempty"`
	URL           string `bson:"url,omitempty" json:"url,omitempty"`
}
