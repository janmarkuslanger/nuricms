package setup

import (
	"fmt"

	"github.com/janmarkuslanger/nuricms/internal/env"
)

func LoadEnv(s env.EnvSource) (*env.Env, error) {
	var e env.Env

	secret := s.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET must be set")
	}
	e.Secret = secret

	return &e, nil
}
