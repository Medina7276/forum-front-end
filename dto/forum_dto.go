package dto

import "forum/model"

type Forum struct {
	User    *model.User
	Posts   []*PostDto
	IsLogin bool
}
