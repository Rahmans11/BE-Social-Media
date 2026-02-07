package model

type Profile struct {
	Id          int     `db:"id"`
	AccountId   int     `db:"account_id"`
	FirstName   *string `db:"first_name"`
	LastName    *string `db:"last_name"`
	Email       string  `db:"email"`
	PhoneNumber *string `db:"phone_number"`
	Avatar      *string `db:"avatar"`
	Bio         *string `db:"bio"`
}
