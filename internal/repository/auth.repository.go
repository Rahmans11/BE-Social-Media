package repository

import (
	"context"
	"log"

	"github.com/Rahmans11/final-phase-3/internal/dto"
	"github.com/Rahmans11/final-phase-3/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepo interface {
	FindByEmail(c context.Context, db DBTX, email string) (model.Auth, error)
	FindByEmailAndPassword(c context.Context, db DBTX, email, password string) (model.Auth, error)
	CheckExistingEmail(c context.Context, db DBTX, email string) (bool, error)
	InsertToUsers(c context.Context, db DBTX, newUser dto.AuthRequest) (model.Auth, error)
	CreateProfile(c context.Context, db DBTX, accountId int, email string) (pgconn.CommandTag, error)
}
type AuthRepository struct {
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (a AuthRepository) FindByEmail(c context.Context, db DBTX, email string) (model.Auth, error) {
	authData := model.Auth{}

	sqlStr := "SELECT id, email, password, role FROM accounts WHERE email = $1"

	err := db.QueryRow(c, sqlStr, email).Scan(&authData.Id, &authData.Email, &authData.Password, &authData.Role)
	if err != nil {
		log.Println(err.Error())
		return model.Auth{}, err
	}

	return authData, nil
}

func (a AuthRepository) FindByEmailAndPassword(c context.Context, db DBTX, email, password string) (model.Auth, error) {
	authData := model.Auth{}

	sqlStr := "SELECT id, email, password, role FROM accounts WHERE email = $1 AND password = $2"

	err := db.QueryRow(c, sqlStr, email, password).Scan(&authData.Id, &authData.Email, &authData.Password, &authData.Role)
	if err != nil {
		log.Println(err.Error())
		return model.Auth{}, err
	}

	return authData, nil
}

func (a AuthRepository) CheckExistingEmail(c context.Context, db DBTX, email string) (bool, error) {
	authData := model.Auth{}

	sqlStr := "SELECT id, email, password, role FROM accounts WHERE email = $1"

	err := db.QueryRow(c, sqlStr, email).Scan(&authData.Id, &authData.Role, &authData.Email, &authData.Password)

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (a AuthRepository) InsertToUsers(c context.Context, db DBTX, newUser dto.AuthRequest) (model.Auth, error) {

	authData := model.Auth{}

	sqlStr := "INSERT INTO accounts (email, password) VALUES ($1, $2) RETURNING id, email, password, role"

	values := []any{newUser.Email, newUser.Password}

	row := db.QueryRow(c, sqlStr, values...)

	if err := row.Scan(&authData.Id, &authData.Email, &authData.Password, &authData.Role); err != nil {
		return model.Auth{}, err
	}

	return authData, nil
}

func (a AuthRepository) CreateProfile(c context.Context, db DBTX, accountId int, email string) (pgconn.CommandTag, error) {

	//authData := model.CreateProfile{}

	sqlStr := "INSERT INTO users (account_id, email) VALUES ($1, $2)"

	values := []any{accountId, email}

	//row := db.QueryRow(c, sqlStr, values...)

	return db.Exec(c, sqlStr, values...)

	// if err := row.Scan(&authData.AccountId, &authData.Email); err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// return nil
}
