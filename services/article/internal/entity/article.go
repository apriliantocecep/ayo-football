package entity

type Article struct {
	Content string `json:"content"`
	UserId  string `json:"user_id"`
	Status  string `json:"status"`
}
