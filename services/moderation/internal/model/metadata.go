package model

type MetadataRequest struct {
	ArticleId string
	Title     string
	Author    string
}

type MetadataResponse struct {
	MetadataId string
}
