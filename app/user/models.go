package user

import "time"

type User struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Session struct {
	Token     string    `db:"token"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
