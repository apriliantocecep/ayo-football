package model

type ArticleRequest struct {
	Content string
	UserId  string
}

type ArticleResponse struct {
	ArticleId string
	Status    string
}
