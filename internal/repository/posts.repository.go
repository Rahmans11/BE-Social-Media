package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Rahmans11/final-phase-3/internal/dto"
	"github.com/Rahmans11/final-phase-3/internal/model"
)

type PostsRepo interface {
	CreatePosts(ctx context.Context, db DBTX, data dto.CreatePosts, id int) (model.Posts, error)
	GetUserId(ctx context.Context, db DBTX, id int) (int, error)
}

type PostsRepository struct {
}

func NewPostsRepository() *PostsRepository {
	return &PostsRepository{}
}

func (p PostsRepository) GetUserId(ctx context.Context, db DBTX, id int) (int, error) {

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

func (p PostsRepository) CreatePosts(ctx context.Context, db DBTX, data dto.CreatePosts, id int) (model.Posts, error) {

	post := model.Posts{}

	var columns []string
	var values []string
	var args []any

	i := 1

	columns = append(columns, "user_id")
	values = append(values, fmt.Sprintf("$%d", i))
	log.Println(id)
	log.Println(1)
	args = append(args, id)
	i++

	if data.Caption != "" {
		columns = append(columns, "caption")
		values = append(values, fmt.Sprintf("$%d", i))
		args = append(args, data.Caption)
		i++
	}

	if data.Image != nil {
		columns = append(columns, "image")
		values = append(values, fmt.Sprintf("$%d", i))
		args = append(args, data.Image.Filename)
		i++
	}

	var sqlStr strings.Builder

	sqlStr.WriteString(fmt.Sprintf("INSERT INTO posts (%s) VALUES (%s) RETURNING id, user_id, caption, image;",
		strings.Join(columns, ", "),
		strings.Join(values, ", ")))

	err := db.QueryRow(ctx, sqlStr.String(), args...).Scan(&post.Id, &post.UserId, &post.Caption, &post.Image)
	if err != nil {
		log.Println(err.Error())
		return model.Posts{}, err
	}

	return post, nil
}
