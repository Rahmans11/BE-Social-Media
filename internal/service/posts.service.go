package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Rahmans11/final-phase-3/internal/dto"
	"github.com/Rahmans11/final-phase-3/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type PostsService struct {
	postsRepository repository.PostsRepo
	redis           *redis.Client
	db              *pgxpool.Pool
}

func NewPostsService(postsRepository repository.PostsRepo, db *pgxpool.Pool, rdb *redis.Client) *PostsService {
	return &PostsService{
		postsRepository: postsRepository,
		redis:           rdb,
		db:              db,
	}
}

func (p *PostsService) CreatePosts(ctx context.Context, data dto.CreatePosts, id int) (dto.Posts, error) {
	tx, etx := p.db.Begin(ctx)
	if etx != nil {
		return dto.Posts{}, fmt.Errorf("failed to begin transaction: %w", etx)
	}

	defer tx.Rollback(ctx)

	userId, err := p.postsRepository.GetUserId(ctx, p.db, id)
	if err != nil {
		return dto.Posts{}, err
	}

	post, e := p.postsRepository.CreatePosts(ctx, p.db, data, userId)
	if e != nil {
		return dto.Posts{}, e
	}

	if etx = tx.Commit(ctx); etx != nil {
		log.Println("failed to commit transaction: %w", etx)
		return dto.Posts{}, fmt.Errorf("failed to commit transaction: %w", etx)
	}

	response := dto.Posts{
		Id:      post.Id,
		UserId:  post.UserId,
		Caption: post.Caption,
		Image:   post.Image,
	}

	rkey := fmt.Sprintf("rahman:social-media:post:%d", response.Id)

	cachestr, e := json.Marshal(response)
	if e != nil {
		log.Println(e.Error())
		log.Println("failed to marshal")
	} else {
		status := p.redis.Set(ctx, rkey, string(cachestr), time.Hour*1)
		if status.Err() != nil {
			log.Println("caching failed")
			log.Println(status.Err().Error())
		}
	}

	return response, nil
}
