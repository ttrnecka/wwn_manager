package service

import (
	"context"

	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/internal/repository"
)

type UserService interface {
	GetByName(context.Context, string) (*entity.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{r}
}

func (s *userService) GetByName(ctx context.Context, name string) (*entity.User, error) {
	return s.repo.GetByField(ctx, "username", name)
}
