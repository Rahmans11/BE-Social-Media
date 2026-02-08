package dto

type AddFollowed struct {
	FollowedId int `json:"followed_id" binding:"required"`
}
