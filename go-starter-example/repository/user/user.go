package user

import (
	"github.com/wzzfarewell/go-mod/infrastructure/database"
	"github.com/wzzfarewell/go-starter-example/domain/user"
	"github.com/wzzfarewell/go-starter-example/repository/user/model"
	"github.com/wzzfarewell/go-starter-example/repository/user/query"
	"gorm.io/gorm"
)

type UserRepository struct {
	q *query.Query
	database.ModelConverter[model.User, user.User]
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepository{
		q: query.Use(db),
	}
}

func (r *UserRepository) Get(id int32) (*user.User, error) {
	m, err := r.q.User.Where(r.q.User.ID.Eq(id)).First()
	return r.ToDomainWithError(m, err)
}

func (r *UserRepository) Create(entity ...*user.User) error {
	models, err := r.ToModels(entity)
	if err != nil {
		return err
	}
	return r.q.User.Create(models...)
}

func (r *UserRepository) Update(entity *user.User) error {
	_, err := r.q.User.Updates(entity)
	return err
}

func (r *UserRepository) Delete(id ...int32) (int64, error) {
	info, err := r.q.User.Where(r.q.User.ID.In(id...)).Delete()
	return info.RowsAffected, err
}

func (r *UserRepository) FindAll() ([]*user.User, error) {
	models, err := r.q.User.Find()
	if err != nil {
		return nil, err
	}
	return r.ToDomains(models)
}
