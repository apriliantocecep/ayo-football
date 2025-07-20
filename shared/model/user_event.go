package model

type UserEvent struct {
	UserId   string `json:"user_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (u *UserEvent) GetId() string {
	return u.UserId
}
