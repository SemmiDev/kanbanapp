package repository

import (
	"context"

	"github.com/SemmiDev/kanbanapp/internal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	user := entity.User{}
	err := r.db.WithContext(ctx).First(&user, id).Error
	return user, err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	user := entity.User{}
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, nil
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
	return err
}
