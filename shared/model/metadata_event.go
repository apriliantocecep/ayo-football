package model

type MetadataEvent struct {
	ArticleId        string `json:"article_id,omitempty"`
	ModerationStatus string `json:"moderation_status,omitempty"`
}

func (u *MetadataEvent) GetId() string {
	return u.ArticleId
}
