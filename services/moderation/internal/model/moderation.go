package model

type CheckContentRequest struct {
	Content string
}

type CheckContentResponse struct {
	IsPass bool
}

type ProcessRequest struct {
	ArticleId string
	IsPass    bool
}

type ProcessResponse struct {
	ModerationStatus string
}
