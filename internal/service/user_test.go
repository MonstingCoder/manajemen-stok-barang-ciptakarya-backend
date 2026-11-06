package service

import (
	"context"
	"testing"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type userRepoStub struct {
	updated *model.User
	foundID uint
}

func (*userRepoStub) Create(context.Context, *model.User) error { return nil }
func (*userRepoStub) Finds(context.Context, uint, uint) ([]model.User, error) {
	return nil, nil
}
func (*userRepoStub) FindByEmail(context.Context, string) (*model.User, error) {
	return nil, nil
}
func (r *userRepoStub) FindByID(_ context.Context, id uint) (*model.User, error) {
	r.foundID = id
	return &model.User{}, nil
}
func (r *userRepoStub) Update(_ context.Context, _ uint, user *model.User) error {
	r.updated = user
	return nil
}
func (*userRepoStub) Delete(context.Context, uint) error { return nil }

func TestUserUpdateAllowsPartialFields(t *testing.T) {
	repo := new(userRepoStub)
	service := NewUserService(repo)
	name := "Kru"

	if err := service.Update(context.Background(), 1, &model.UserIn{Name: &name}); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if repo.updated == nil || repo.updated.Name == nil || *repo.updated.Name != name {
		t.Fatal("updated name was not passed to the repository")
	}
	if repo.updated.Email != nil || repo.updated.PasswordHash != nil || repo.updated.RoleID != nil {
		t.Fatal("unspecified fields should not be included in the partial update")
	}
}

func TestUserUpdateRejectsEmptyFields(t *testing.T) {
	repo := new(userRepoStub)
	service := NewUserService(repo)

	if err := service.Update(context.Background(), 1, new(model.UserIn)); err == nil {
		t.Fatal("Update with no fields should return an error")
	}

	if repo.updated != nil {
		t.Fatal("repository should not be called for an empty update")
	}
}

func TestUserFindByID(t *testing.T) {
	repo := new(userRepoStub)
	service := NewUserService(repo)

	if _, err := service.FindByID(context.Background(), 42); err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if repo.foundID != 42 {
		t.Fatalf("FindByID called with %d, want 42", repo.foundID)
	}
}
