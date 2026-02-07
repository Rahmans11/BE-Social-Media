package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Rahmans11/final-phase-3/internal/dto"
	"github.com/Rahmans11/final-phase-3/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProfileRepo interface {
	FindProfile(c context.Context, db DBTX, id int) (model.Profile, error)
	EditProfile(ctx context.Context, db DBTX, data dto.EditProfile, id int) (pgconn.CommandTag, error)
	FindAvatar(ctx context.Context, db DBTX, id int) (string, error)
	FindOldPassword(ctx context.Context, db DBTX, id int) (string, error)
	ChangePassword(ctx context.Context, db DBTX, newPassword string, id int) (pgconn.CommandTag, error)
}

type ProfileRepository struct {
}

func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{}
}

func (p ProfileRepository) FindProfile(c context.Context, db DBTX, accountId int) (model.Profile, error) {
	profile := model.Profile{}

	sqlStr := `
	SELECT id, account_id, first_name, last_name, email, phone_number, avatar, bio
	FROM users 
	WHERE account_id = $1
	`
	err := db.QueryRow(c, sqlStr, accountId).Scan(&profile.Id, &profile.AccountId, &profile.FirstName, &profile.LastName,
		&profile.Email, &profile.PhoneNumber, &profile.Avatar, &profile.Bio)
	if err != nil {
		log.Println(err.Error())
		return model.Profile{}, err
	}

	return profile, nil
}

func (p ProfileRepository) EditProfile(ctx context.Context, db DBTX, data dto.EditProfile, id int) (pgconn.CommandTag, error) {

	var args []any
	paramIndex := 1
	fieldCount := 0

	var sqlStr strings.Builder
	sqlStr.WriteString(`
	UPDATE users SET
	`)

	if data.FirstName != nil {
		if fieldCount > 0 {
			sqlStr.WriteString(", ")
		}
		sqlStr.WriteString(fmt.Sprintf("first_name = $%d", paramIndex))
		args = append(args, data.FirstName)
		paramIndex++
		fieldCount++
	}

	if data.LastName != nil {
		if fieldCount > 0 {
			sqlStr.WriteString(", ")
		}
		sqlStr.WriteString(fmt.Sprintf("last_name = $%d", paramIndex))
		args = append(args, data.LastName)
		paramIndex++
		fieldCount++
	}

	if data.PhoneNumber != nil {
		if fieldCount > 0 {
			sqlStr.WriteString(", ")
		}
		sqlStr.WriteString(fmt.Sprintf("phone_number = $%d", paramIndex))
		args = append(args, data.PhoneNumber)
		paramIndex++
		fieldCount++
	}

	if data.Avatar != nil {
		if fieldCount > 0 {
			sqlStr.WriteString(", ")
		}
		sqlStr.WriteString(fmt.Sprintf("avatar = $%d", paramIndex))
		args = append(args, data.Avatar.Filename)
		paramIndex++
		fieldCount++
	}

	if data.Bio != nil {
		if fieldCount > 0 {
			sqlStr.WriteString(", ")
		}
		sqlStr.WriteString(fmt.Sprintf("bio = $%d", paramIndex))
		args = append(args, data.Bio)
		paramIndex++
		fieldCount++
	}

	sqlStr.WriteString(fmt.Sprintf(" WHERE account_id = $%d", paramIndex))

	args = append(args, id)

	log.Println(sqlStr.String())

	return db.Exec(ctx, sqlStr.String(), args...)
}

func (p ProfileRepository) FindAvatar(ctx context.Context, db DBTX, id int) (string, error) {

	var photo sql.NullString

	sqlStr := `
	SELECT avatar
	FROM users 
	WHERE account_id = $1
	`
	err := db.QueryRow(ctx, sqlStr, id).Scan(&photo)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return photo.String, nil
}

func (p ProfileRepository) FindOldPassword(ctx context.Context, db DBTX, id int) (string, error) {

	var oldPassword string

	sqlStr := `
	SELECT password, id
	FROM users 
	WHERE id = $1
	`
	err := db.QueryRow(ctx, sqlStr, id).Scan(&oldPassword, &id)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return oldPassword, nil
}

func (p ProfileRepository) ChangePassword(ctx context.Context, db DBTX, newPassword string, id int) (pgconn.CommandTag, error) {

	sqlStr := `
	UPDATE users SET password = $1
	WHERE id = $2
	`

	return db.Exec(ctx, sqlStr, newPassword, id)
}
