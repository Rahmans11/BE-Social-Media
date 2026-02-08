package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Rahmans11/final-phase-3/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type FollowsService struct {
	followsRepository repository.FollowsRepo
	redis             *redis.Client
	db                *pgxpool.Pool
}

func NewFollowsService(followsRepository repository.FollowsRepo, db *pgxpool.Pool, rdb *redis.Client) *FollowsService {
	return &FollowsService{
		followsRepository: followsRepository,
		redis:             rdb,
		db:                db,
	}
}

func (f FollowsService) AddFollowed(ctx context.Context, id, followedId int) error {

	tx, etx := f.db.Begin(ctx)
	if etx != nil {
		return fmt.Errorf("failed to begin transaction: %w", etx)
	}

	defer tx.Rollback(ctx)

	followerId, err := f.followsRepository.GetUserId(ctx, f.db, id)
	if err != nil {
		return err
	}
	log.Println(followerId)
	log.Println("test")

	cmd, e := f.followsRepository.AddFollowed(ctx, f.db, followerId, followedId)
	if e != nil {
		return e
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("No data updated")
	}

	if etx = tx.Commit(ctx); etx != nil {
		log.Println("failed to commit transaction: %w", etx)
		return fmt.Errorf("failed to commit transaction: %w", etx)
	}

	return nil
}
