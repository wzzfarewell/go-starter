package user

import (
	"github.com/wzzfarewell/go-mod/infrastructure/database"
	"github.com/wzzfarewell/go-starter-example/domain/user"
	"github.com/wzzfarewell/go-starter-example/repository/user/model"
	"github.com/wzzfarewell/go-starter-example/repository/user/query"
	"gorm.io/gorm"
)

type UserInfoRepository struct {
	q *query.Query
	database.ModelConverter[model.UserInfo, user.UserInfo]
}

func NewUserInfoRepository(db *gorm.DB) user.UserInfoRepository {
	return &UserInfoRepository{
		q: query.Use(db),
	}
}

func (r *UserInfoRepository) Get(id int32) (*user.UserInfo, error) {
	m, err := r.q.UserInfo.Where(r.q.UserInfo.ID.Eq(id)).First()
	return r.ToDomainWithError(m, err)
}

func (r *UserInfoRepository) Create(entity ...*user.UserInfo) error {
	models, err := r.ToModels(entity)
	if err != nil {
		return err
	}
	return r.q.UserInfo.Create(models...)
}

func (r *UserInfoRepository) Update(entity *user.UserInfo) error {
	_, err := r.q.UserInfo.Updates(entity)
	return err
}

func (r *UserInfoRepository) Delete(id ...int32) (int64, error) {
	info, err := r.q.UserInfo.Where(r.q.UserInfo.ID.In(id...)).Delete()
	return info.RowsAffected, err
}

func (r *UserInfoRepository) FindAll() ([]*user.UserInfo, error) {
	models, err := r.q.UserInfo.Find()
	if err != nil {
		return nil, err
	}
	return r.ToDomains(models)
}
