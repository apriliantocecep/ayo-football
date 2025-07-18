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

type ModerationRequest struct {
	ArticleId string
	UserId    string
}

type ArticleResource struct {
	ID     string
	UserId string
}

type PublishArticleRequest struct {
	ArticleId        string
	ModerationStatus string
}

type PublishArticleResponse struct {
	Status string
}
