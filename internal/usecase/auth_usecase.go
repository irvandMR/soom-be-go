package usecase

import (
	"errors"
	"os"
	"soom-be-go/internal/domain"
	"soom-be-go/internal/middleware"
	"soom-be-go/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	repo repository.AuthRepository
}

func NewAuthUsecase(r repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{repo: r}
}

func (u *AuthUsecase) Register(req domain.RegisterRequest) error {

	// Check Email Exist
	_, err := u.repo.FindByEmail(req.Email)
	if err == nil {
		return &middleware.GlobalError{
			Code:    "EMAIL_ALREADY_EXISTS",
			Message: "Email already exists",
			Status:  400,
		}
	}

	// hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return middleware.ErrInternalServer
	}

	user := &domain.User{
		BaseModelWithDeleted: domain.BaseModelWithDeleted{
			BaseModel: domain.BaseModel{
				CreatedBy: req.Username,
			},
		},
		TenantId:   req.TenantId,
		Username:   req.Username,
		Email:      req.Email,
		Password:   string(hashPassword),
		TenantRole: "Employee",
		Role:       "user",
	}

	return u.repo.Register(user)
}

func (u *AuthUsecase) Login(req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &middleware.GlobalError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Email or password is incorrect",
				Status:  401,
			}
		}
		return nil, middleware.ErrInternalServer
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, &middleware.GlobalError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Email or password is incorrect",
			Status:  401,
		}
	}

	// generate access token — expired 15 menit
	accessToken, err := generateToken(user, 15*time.Minute)
	if err != nil {
		return nil, middleware.ErrInternalServer
	}

	// generate refresh token — expired 7 hari
	refreshToken, err := generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, middleware.ErrInternalServer
	}

	// simpan refresh token ke DB
	err = u.repo.SaveRefreshToken(&domain.RefreshToken{
		BaseModel: domain.BaseModel{
			CreatedBy: req.Email,
		},
		UserId:    user.Id,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})

	if err != nil {
		return nil, middleware.ErrInternalServer
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
	}, nil
}

func (u *AuthUsecase) RefreshToken(req domain.RefreshTokenRequest) (*domain.LoginResponse, error) {
	// cek refresh token di DB
	savedToken, err := u.repo.FindRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, &middleware.GlobalError{
			Code:    "INVALID_REFRESH_TOKEN",
			Message: "Refresh token is invalid or expired",
			Status:  401,
		}
	}

	// ambil data user
	user, err := u.repo.FindById(savedToken.UserId)
	if err != nil {
		return nil, middleware.ErrInternalServer
	}

	// generate access token baru
	accessToken, err := generateToken(user, 15*time.Minute)
	if err != nil {
		return nil, middleware.ErrInternalServer
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
	}, nil
}

func (u *AuthUsecase) Logout(refreshToken string) error {
	return u.repo.DeleteRefreshToken(refreshToken)
}

func generateToken(user *domain.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
