package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/wzzfarewell/go-mod/infrastructure/utils/copyutil"
	"github.com/wzzfarewell/go-starter-example/domain/user"
)

type UserService struct {
	UserRepo user.UserRepository
}

func NewUserService(UserRepo user.UserRepository) user.UserService {
	return &UserService{
		UserRepo: UserRepo,
	}
}

func (s *UserService) Get(ctx context.Context, id int32) (*user.User, error) {
	return s.UserRepo.Get(id)
}

func (s *UserService) Create(ctx context.Context, entity ...*user.User) error {
	return s.UserRepo.Create(entity...)
}

func (s *UserService) Update(ctx context.Context, entity *user.User) (*user.User, error) {
	m, err := s.UserRepo.Get(entity.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get model from repository")
	}
	if m == nil {
		return nil, errors.New("model not found")
	}
	if err := copyutil.CopyTo(entity, m); err != nil {
		return nil, errors.Wrap(err, "failed to copy entity to model")
	}
	if err := s.UserRepo.Update(m); err != nil {
		return nil, errors.Wrap(err, "failed to update model in repository")
	}
	return m, nil
}

func (s *UserService) Delete(ctx context.Context, id ...int32) error {
	_, err := s.UserRepo.Delete(id...)
	return err
}

func (s *UserService) FindAll(ctx context.Context) ([]*user.User, error) {
	return s.UserRepo.FindAll()
}
