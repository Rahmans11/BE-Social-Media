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

type ProfileService struct {
	profileRepository repository.ProfileRepo
	redis             *redis.Client
	db                *pgxpool.Pool
}

func NewProfileService(profileRepository repository.ProfileRepo, db *pgxpool.Pool, rdb *redis.Client) *ProfileService {
	return &ProfileService{
		profileRepository: profileRepository,
		redis:             rdb,
		db:                db,
	}
}

func (p *ProfileService) GetProfile(ctx context.Context, id int) (dto.ProfileResponse, error) {

	rkey := fmt.Sprintf("rahman:social-media:profile:%d", id)
	rsc := p.redis.Get(ctx, rkey)

	if rsc.Err() == nil {
		var result dto.ProfileResponse
		cache, e := rsc.Bytes()
		if e != nil {
			log.Println(e.Error())
		} else {
			e := json.Unmarshal(cache, &result)
			if e != nil {
				log.Println(e.Error())
			} else {
				return result, nil
			}
		}
	}

	if rsc.Err() == redis.Nil {
		log.Println("profile cache miss")
	}

	data, e := p.profileRepository.FindProfile(ctx, p.db, id)
	if e != nil {
		return dto.ProfileResponse{}, e
	}

	response := dto.ProfileResponse{
		UserId:      data.Id,
		AccountId:   data.AccountId,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Avatar:      data.Avatar,
		Bio:         data.Bio,
	}

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

func (p *ProfileService) GetOtherProfile(ctx context.Context, id int) (dto.ProfileOtherResponse, error) {

	rkey := fmt.Sprintf("rahman:social-media:profile:user:%d", id)
	rsc := p.redis.Get(ctx, rkey)

	if rsc.Err() == nil {
		var result dto.ProfileOtherResponse
		cache, e := rsc.Bytes()
		if e != nil {
			log.Println(e.Error())
		} else {
			e := json.Unmarshal(cache, &result)
			if e != nil {
				log.Println(e.Error())
			} else {
				return result, nil
			}
		}
	}

	if rsc.Err() == redis.Nil {
		log.Println("profile cache miss")
	}

	data, e := p.profileRepository.FindProfile(ctx, p.db, id)
	if e != nil {
		return dto.ProfileOtherResponse{}, e
	}

	log.Println(data)

	response := dto.ProfileOtherResponse{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Avatar:    data.Avatar,
		Bio:       data.Bio,
	}

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

func (p *ProfileService) UpdateProfile(ctx context.Context, data dto.EditProfile, id int) error {

	cmd, e := p.profileRepository.EditProfile(ctx, p.db, data, id)
	if e != nil {
		return e
	}

	rkey := fmt.Sprintf("rahman:social-media:profile:%d", id)
	deleted, err := p.redis.Del(ctx, rkey).Result()
	if err != nil {
		return err
	}

	if deleted == 0 {
		fmt.Printf("Key %s not found in Redis (or already deleted)\n", rkey)
	} else {
		fmt.Printf("Key %s deleted successfully\n", rkey)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("No data updated")
	}

	return nil
}

func (p *ProfileService) GetPhoto(ctx context.Context, id int) (string, error) {

	data, err := p.profileRepository.FindAvatar(ctx, p.db, id)
	if err != nil {
		return "", err
	}

	photo := data

	return photo, nil
}

// func (p *ProfileService) ChangePassword(ctx context.Context, data dto.ChangePassword, id int) error {

// 	tx, e := p.db.Begin(ctx)
// 	if e != nil {
// 		return fmt.Errorf("failed to begin transaction: %w", e)
// 	}

// 	defer tx.Rollback(ctx)

// 	op, e := p.profileRepository.FindOldPassword(ctx, tx, id)
// 	log.Println(data.OldPassword)
// 	if e != nil {
// 		return fmt.Errorf("no user fount")
// 	}

// 	hc := pkg.HashConfig{}

// 	pw, e := hc.ComparePwdAndHash(data.OldPassword, op)
// 	if e != nil {
// 		return e
// 	}

// 	if !pw {
// 		return err.WrongPassword
// 	}

// 	if len(data.NewPassword) < 6 {
// 		return err.InvalidFormatPassword
// 	}

// 	var hasUpper, hasLower, hasDigit, hasSpecial bool

// 	for _, char := range data.NewPassword {
// 		if unicode.IsUpper(char) {
// 			hasUpper = true
// 		} else if unicode.IsLower(char) {
// 			hasLower = true
// 		} else if unicode.IsDigit(char) {
// 			hasDigit = true
// 		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
// 			hasSpecial = true
// 		}
// 	}

// 	if !hasUpper {
// 		return err.InvalidFormatPassword
// 	}
// 	if !hasLower {
// 		return err.InvalidFormatPassword
// 	}
// 	if !hasDigit {
// 		return err.InvalidFormatPassword
// 	}
// 	if !hasSpecial {
// 		return err.InvalidFormatPassword
// 	}

// 	hc.UseRecommended()

// 	hp, err := hc.GenHash(data.NewPassword)
// 	if err != nil {
// 		return err
// 	}

// 	data.NewPassword = hp

// 	cmd, e := p.profileRepository.ChangePassword(ctx, tx, data.NewPassword, id)
// 	log.Println(data.NewPassword)
// 	if e != nil {
// 		return fmt.Errorf("failed to change password")
// 	}

// 	if cmd.RowsAffected() == 0 {
// 		return fmt.Errorf("No data updated")
// 	}

// 	if err = tx.Commit(ctx); err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}

// 	return nil
// }
