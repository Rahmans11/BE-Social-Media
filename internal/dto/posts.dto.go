package dto

import "mime/multipart"

type CreatePosts struct {
	Caption string                `form:"caption" example:"test caption" binding:"required"`
	Image   *multipart.FileHeader `form:"image"`
}

type Posts struct {
	Id      int     `json:"id"`
	UserId  int     `json:"user_id"`
	Caption string  `json:"caption"`
	Image   *string `json:"image"`
}
