package user

import (
	"context"
	"time"
)

type User struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(entity ...*User) error
	Update(entity *User) error
	Delete(id ...int32) (int64, error)
	Get(id int32) (*User, error)
	FindAll() ([]*User, error)
}

type UserService interface {
	Create(ctx context.Context, entity ...*User) error
	Update(ctx context.Context, entity *User) (*User, error)
	Delete(ctx context.Context, id ...int32) error
	Get(ctx context.Context, id int32) (*User, error)
	FindAll(ctx context.Context) ([]*User, error)
}
