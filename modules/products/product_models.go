package products

import (
	"my-pp/modules/users"
	"time"
)

type Body struct {
	Title    string      `json:"title"`
	Text     string      `json:"text"`
	UserId   string      `json:"userId"`
	UpdateAt time.Time   `json:"updateAt"`
	User     *users.User `json:"user,omitempty"`
}

type ItemData struct {
	ID   string `json:"id"`
	Body Body   `json:"body"`
}
