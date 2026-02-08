package model

type Follows struct {
	FollowerId int `db:"follower_id"`
	FollowedId int `db:"followed_id"`
}
