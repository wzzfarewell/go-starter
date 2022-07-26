package user

import (
	"context"
)

type UserInfo struct {
	ID       int32  `json:"id"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type UserInfoRepository interface {
	Create(entity ...*UserInfo) error
	Update(entity *UserInfo) error
	Delete(id ...int32) (int64, error)
	Get(id int32) (*UserInfo, error)
	FindAll() ([]*UserInfo, error)
}

type UserInfoService interface {
	Create(ctx context.Context, entity ...*UserInfo) error
	Update(ctx context.Context, entity *UserInfo) (*UserInfo, error)
	Delete(ctx context.Context, id ...int32) error
	Get(ctx context.Context, id int32) (*UserInfo, error)
	FindAll(ctx context.Context) ([]*UserInfo, error)
}
