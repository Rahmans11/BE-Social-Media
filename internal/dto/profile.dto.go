package dto

import "mime/multipart"

type ProfileParam struct {
	Id int `uri:"id"`
}

type ProfileResponse struct {
	UserId      int     `json:"user_id"`
	AccountId   int     `json:"account_id"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Avatar      *string `json:"avatar"`
	Bio         *string `json:"bio"`
}

type ProfileOtherResponse struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     string  `json:"email"`
	Avatar    *string `json:"avatar"`
	Bio       *string `json:"bio"`
}

type EditProfile struct {
	FirstName   *string               `form:"first_name"`
	LastName    *string               `form:"last_name"`
	PhoneNumber *string               `form:"phone_number"`
	Avatar      *multipart.FileHeader `form:"avatar"`
	Bio         *string               `form:"bio"`
}
