package repository

import (
	"soom-be-go/internal/domain"
	"time"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(email string) (*domain.User, error)
	Register(user *domain.User) error
	SaveRefreshToken(token *domain.RefreshToken) error
	FindRefreshToken(token string) (*domain.RefreshToken, error)
	DeleteRefreshToken(token string) error
	FindById(id string) (*domain.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ? and deleted_at is null", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *authRepository) Register(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) SaveRefreshToken(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *authRepository) FindRefreshToken(token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	result := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &refreshToken, nil
}

func (r *authRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&domain.RefreshToken{}).Error
}

func (r *authRepository) FindById(id string) (*domain.User, error) {
    var user domain.User
    result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
