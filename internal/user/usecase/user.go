package usecase

import (
	"context"

	"project-evredika/internal/user"
	v1 "project-evredika/pkg/api/v1"
)

type userUsecase struct {
	storage storage
}

func (u *userUsecase) CreateUser(ctx context.Context, user *v1.User) (err error) {
	return u.storage.CreateUser(ctx, user)
}
func (u *userUsecase) DeleteUser(ctx context.Context, ID string) (err error) {
	return u.storage.DeleteUser(ctx, ID)
}
func (u *userUsecase) UpdateUser(ctx context.Context, user *v1.User) (err error) {
	return u.storage.UpdateUser(ctx, user)
}
func (u *userUsecase) GetUser(ctx context.Context, ID string) (user *v1.User, err error) {
	return u.storage.GetUser(ctx, ID)
}
func (u *userUsecase) ListUsers(ctx context.Context, skip, limit int) (users []*v1.User) {
	return u.storage.ListUsers(ctx, skip, limit)
}

// NewUserUsecase ..
func NewUserUsecase(storage storage) user.Usecase {
	return &userUsecase{storage: storage}
}
