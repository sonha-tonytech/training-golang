package types

import "time"

type User struct {
	ID   string `json:"id"`
	Body struct {
		Name     string `json:"name"`
		UserName string `json:"userName"`
		Password string `json:"password"`
	} `json:"body"`
}

type Body struct {
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	UserId   string    `json:"userId"`
	UpdateAt time.Time `json:"updateAt"`
	User     *User     `json:"user,omitempty"`
}

type ItemData struct {
	ID   string `json:"id"`
	Body Body   `json:"body"`
}
