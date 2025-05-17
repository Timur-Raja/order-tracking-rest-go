package usersql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/timur-raja/order-tracking-rest-go/app/user"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/timur-raja/order-tracking-rest-go/db"
)

type insertUserQuery struct {
	db.BaseQuery
	Values struct {
		User user.User
	}
	Returned struct {
		ID int `db:"id"`
	}
}

func NewInsertUserQuery(conn *pgxpool.Pool, u user.User) *insertUserQuery {
	return &insertUserQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		Values:    struct{ User user.User }{User: u},
	}
}

func (q *insertUserQuery) Run(ctx context.Context) error {
	query := `
        INSERT INTO users (email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
    `
	// Execute the query and scan the returned ID into q.Returned.ID
	err := pgxscan.Get(ctx, q.DBConn, &q.Returned, query,
		q.Values.User.Email,
		q.Values.User.Password,
		q.Values.User.Name,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// the user view will be the user object that we return to the FE, including eventual necessary extra fields from other tables
type selectUserViewByIDQuery struct {
	db.BaseQuery
	Where struct {
		ID int
	}
	UserView *user.UserView
}

func NewSelectUserViewByIDQuery(Conn *pgxpool.Pool, id int) *selectUserViewByIDQuery {
	return &selectUserViewByIDQuery{
		BaseQuery: db.BaseQuery{DBConn: Conn},
		Where:     struct{ ID int }{ID: id},
		UserView:  new(user.UserView),
	}
}

func (q *selectUserViewByIDQuery) Run(ctx context.Context) error {
	q.UserView = new(user.UserView)
	query := `
        SELECT * FROM users_view WHERE id = $1;
    `
	// Execute the query and scan the returned ID into q.Returned.ID
	err := pgxscan.Get(ctx, q.DBConn, q.UserView, query,
		q.Where.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
