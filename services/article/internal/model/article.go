package model

type ArticleRequest struct {
	Title   string
	Author  string
	Content string
	UserId  string
}

type ArticleResponse struct {
	ArticleId string
	Status    string
}
