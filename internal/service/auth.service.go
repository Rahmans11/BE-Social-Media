package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/Rahmans11/final-phase-3/internal/dto"
	"github.com/Rahmans11/final-phase-3/internal/err"
	"github.com/Rahmans11/final-phase-3/internal/repository"
	"github.com/Rahmans11/final-phase-3/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	authRepository repository.AuthRepo
	redis          *redis.Client
	db             *pgxpool.Pool
}

func NewAuthService(authRepository repository.AuthRepo, db *pgxpool.Pool, rdb *redis.Client) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		redis:          rdb,
		db:             db,
	}
}

func (a AuthService) Login(c context.Context, loginData dto.AuthRequest) (dto.AuthResponse, error) {

	data, e := a.authRepository.FindByEmail(c, a.db, loginData.Email)
	if e != nil {
		return dto.AuthResponse{}, e
	}

	hp := data.Password

	if len(hp) == 0 {
		return dto.AuthResponse{}, errors.New("email or password is wrong")
	}
	hc := pkg.HashConfig{}
	p, e := hc.ComparePwdAndHash(loginData.Password, hp)
	if e != nil {
		return dto.AuthResponse{}, e
	}

	if !p {
		return dto.AuthResponse{}, errors.New("email or password is wrong")
	}

	response := dto.AuthResponse{
		Id:    data.Id,
		Email: data.Email,
		Role:  data.Role,
	}

	token, e := a.GenJWTToken(response)
	if e != nil {
		return dto.AuthResponse{}, err.FailedGenerateToken
	}

	response.Token = token

	rkey := fmt.Sprintf("rahman:social-media:whitelist-token:%d", data.Id)

	status := a.redis.Set(c, rkey, token, time.Hour*24)
	if status.Err() != nil {
		log.Println("caching failed")
		log.Println(status.Err().Error())
	}

	return response, nil
}

func (a AuthService) Register(c context.Context, newUser dto.AuthRequest) (dto.AuthResponse, error) {

	tx, etx := a.db.Begin(c)
	if etx != nil {
		return dto.AuthResponse{}, fmt.Errorf("failed to begin transaction: %w", etx)
	}

	defer tx.Rollback(c)

	exists, e := a.authRepository.CheckExistingEmail(c, a.db, newUser.Email)
	if e != nil {

		return dto.AuthResponse{}, e
	}

	if exists {

		return dto.AuthResponse{}, err.ExistingEmail
	}

	email := strings.TrimSpace(newUser.Email)
	if email == "" {
		return dto.AuthResponse{}, e
	}

	if !strings.Contains(email, "@") {
		return dto.AuthResponse{}, err.InvalidFormatEmail
	}

	if strings.HasPrefix(email, "@") || strings.HasSuffix(email, "@") {
		return dto.AuthResponse{}, err.InvalidFormatEmail
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return dto.AuthResponse{}, err.InvalidFormatEmail
	}

	localPart := parts[0]
	domain := parts[1]

	if localPart == "" {
		return dto.AuthResponse{}, err.InvalidFormatEmail
	}

	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return dto.AuthResponse{}, err.InvalidFormatEmail
	}

	if len(newUser.Password) < 6 {
		return dto.AuthResponse{}, err.InvalidFormatPassword
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range newUser.Password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}

	if !hasUpper {
		return dto.AuthResponse{}, err.InvalidFormatPassword
	}
	if !hasLower {
		return dto.AuthResponse{}, err.InvalidFormatPassword
	}
	if !hasDigit {
		return dto.AuthResponse{}, err.InvalidFormatPassword
	}
	if !hasSpecial {
		return dto.AuthResponse{}, err.InvalidFormatPassword
	}

	hc := pkg.HashConfig{}
	hc.UseRecommended()

	hp, e := hc.GenHash(newUser.Password)
	if e != nil {
		return dto.AuthResponse{}, e
	}

	newUser.Password = hp

	data, e := a.authRepository.InsertToUsers(c, a.db, newUser)
	if e != nil {
		return dto.AuthResponse{}, e
	}

	cmd, eCP := a.authRepository.CreateProfile(c, a.db, data.Id, data.Email)
	if eCP != nil {
		log.Println("error jir")
		return dto.AuthResponse{}, eCP
	}

	if cmd.RowsAffected() == 0 {
		return dto.AuthResponse{}, fmt.Errorf("No data updated")
	}

	if etx = tx.Commit(c); etx != nil {
		log.Println("failed to commit transaction: %w", etx)
		return dto.AuthResponse{}, fmt.Errorf("failed to commit transaction: %w", etx)
	}

	response := dto.AuthResponse{
		Id:    data.Id,
		Email: data.Email,
		Role:  data.Role,
	}

	token, e := a.GenJWTToken(response)
	if e != nil {
		return dto.AuthResponse{}, err.FailedGenerateToken
	}

	response.Token = token

	rkey := fmt.Sprintf("rahman:social-media:whitelist-token:%d", data.Id)
	status := a.redis.Set(c, rkey, token, time.Hour*24)
	if status.Err() != nil {
		log.Println("caching failed")
		log.Println(status.Err().Error())
	}

	return response, nil
}

func (a AuthService) GenJWTToken(user dto.AuthResponse) (string, error) {
	claims := pkg.NewJWTClaims(user.Id, user.Role)
	return claims.GenToken()
}

func (a AuthService) Logout(ctx context.Context, id int) error {

	rkey := fmt.Sprintf("rahman:social-media:whitelist-token:%d", id)
	deleted, err := a.redis.Del(ctx, rkey).Result()
	if err != nil {
		return err
	}

	if deleted == 0 {
		fmt.Printf("Key %s not found in Redis (or already deleted)\n", rkey)
	} else {
		fmt.Printf("Key %s deleted successfully\n", rkey)
	}

	return nil
}
