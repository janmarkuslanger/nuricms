package setup

import (
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
)

type userCreator interface {
	Create(dto dto.UserData) (*model.User, error)
	List(page, pageSize int) ([]model.User, int64, error)
}

func InitAdminUser(s userCreator) error {
	_, count, err := s.List(1, 1)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	if _, err = s.Create(dto.UserData{
		Email:    "admin@admin.com",
		Password: "mysecret",
		Role:     string(model.RoleAdmin),
	}); err != nil {
		return err
	}

	return nil
}
