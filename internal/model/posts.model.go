package model

type Posts struct {
	Id      int     `db:"id"`
	UserId  int     `db:"user_id"`
	Caption string  `db:"caption"`
	Image   *string `db:"image"`
}
