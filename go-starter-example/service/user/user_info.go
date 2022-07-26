package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/wzzfarewell/go-mod/infrastructure/utils/copyutil"
	"github.com/wzzfarewell/go-starter-example/domain/user"
)

type UserInfoService struct {
	UserInfoRepo user.UserInfoRepository
}

func NewUserInfoService(UserInfoRepo user.UserInfoRepository) user.UserInfoService {
	return &UserInfoService{
		UserInfoRepo: UserInfoRepo,
	}
}

func (s *UserInfoService) Get(ctx context.Context, id int32) (*user.UserInfo, error) {
	return s.UserInfoRepo.Get(id)
}

func (s *UserInfoService) Create(ctx context.Context, entity ...*user.UserInfo) error {
	return s.UserInfoRepo.Create(entity...)
}

func (s *UserInfoService) Update(ctx context.Context, entity *user.UserInfo) (*user.UserInfo, error) {
	m, err := s.UserInfoRepo.Get(entity.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get model from repository")
	}
	if m == nil {
		return nil, errors.New("model not found")
	}
	if err := copyutil.CopyTo(entity, m); err != nil {
		return nil, errors.Wrap(err, "failed to copy entity to model")
	}
	if err := s.UserInfoRepo.Update(m); err != nil {
		return nil, errors.Wrap(err, "failed to update model in repository")
	}
	return m, nil
}

func (s *UserInfoService) Delete(ctx context.Context, id ...int32) error {
	_, err := s.UserInfoRepo.Delete(id...)
	return err
}

func (s *UserInfoService) FindAll(ctx context.Context) ([]*user.UserInfo, error) {
	return s.UserInfoRepo.FindAll()
}
