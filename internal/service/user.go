package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repos     *repository.Set
	jwtSecret []byte
}

func NewUserService(repos *repository.Set, jwtSecret []byte) *UserService {
	return &UserService{repos: repos, jwtSecret: jwtSecret}
}

func (s *UserService) List(page, pageSize int) ([]model.User, int64, error) {
	return s.repos.User.List(page, pageSize)
}

func (s *UserService) DeleteByID(id uint) error {
	user, err := s.FindByID(id)
	if err != nil {
		return err
	}

	return s.repos.User.Delete(user)
}

func (s *UserService) Create(email, password string, role model.Role) (*model.User, error) {
	switch role {
	case model.RoleAdmin, model.RoleEditor:

	default:
		return nil, errors.New("Not a valid role")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    email,
		Password: string(hash),
		Role:     role,
	}
	if err := s.repos.User.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindByID(id uint) (*model.User, error) {
	return s.repos.User.FindByID(id)
}

func (s *UserService) Save(user *model.User) error {
	return s.repos.User.Save(user)
}

func (s *UserService) Delete(user *model.User) error {
	return s.repos.User.Delete(user)
}

func (s *UserService) LoginUser(email, password string) (string, error) {
	user, err := s.repos.User.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  string(user.Role),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

func (s *UserService) ValidateJWT(tokenStr string) (userID uint, email string, role model.Role, err error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !tok.Valid {
		return 0, "", "", err
	}
	claims := tok.Claims.(jwt.MapClaims)
	uid := uint(claims["sub"].(float64))
	em := claims["email"].(string)
	rl := model.Role(claims["role"].(string))
	return uid, em, rl, nil
}
