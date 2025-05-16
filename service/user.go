package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret []byte
}

func NewUserService(r *repository.UserRepository, jwtSecret []byte) *UserService {
	return &UserService{repo: r, jwtSecret: jwtSecret}
}

func (s *UserService) List() ([]model.User, error) {
	return s.repo.List()
}

func (s *UserService) DeleteByID(id uint) error {
	user, err := s.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(user)
}

func (s *UserService) CreateUser(email, password string, role model.Role) (*model.User, error) {
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
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) Save(user *model.User) error {
	return s.repo.Save(user)
}

func (s *UserService) Delete(user *model.User) error {
	return s.repo.Delete(user)
}

func (s *UserService) LoginUser(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
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
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
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
