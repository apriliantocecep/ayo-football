package model

type ArticleRequest struct {
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content" validate:"required,min=3"`
}

type PublishArticleRequest struct {
	ArticleId string `json:"article_id" validate:"required"`
}
