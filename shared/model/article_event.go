package model

type ArticleEvent struct {
	ID      string `json:"article_id,omitempty"`
	Title   string `json:"title,omitempty"`
	Author  string `json:"author,omitempty"`
	Content string `json:"content,omitempty"`
}

func (u *ArticleEvent) GetId() string {
	return u.ID
}
