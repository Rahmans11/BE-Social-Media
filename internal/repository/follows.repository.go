package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
)

type FollowsRepo interface {
	GetUserId(ctx context.Context, db DBTX, id int) (int, error)
	AddFollowed(ctx context.Context, db DBTX, followerId, followedId int) (pgconn.CommandTag, error)
}

type FollowsRepository struct {
}

func NewFollowsRepository() *FollowsRepository {
	return &FollowsRepository{}
}

func (f FollowsRepository) AddFollowed(ctx context.Context, db DBTX, followerId, followedId int) (pgconn.CommandTag, error) {

	sqlStr := `
	INSERT INTO follows (follower_id, followed_id) values ($1, $2);
	`

	return db.Exec(ctx, sqlStr, followerId, followedId)
}

func (f FollowsRepository) GetUserId(ctx context.Context, db DBTX, id int) (int, error) {

	var userId int

	sqlStr := `
	SELECT id
	FROM users 
	WHERE account_id = $1
	`
	err := db.QueryRow(ctx, sqlStr, id).Scan(&userId)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return userId, nil
}
