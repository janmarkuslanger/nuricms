package setup_test

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/setup"
	"github.com/janmarkuslanger/nuricms/pkg/config"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
)

type SuccessEnv struct{}

func (te SuccessEnv) Getenv(v string) string {
	envs := map[string]string{
		"JWT_SECRET": "hiitsme",
	}

	return envs[v]
}

func TestSetupApp(t *testing.T) {
	dl := sqlite.Open(":memory:")
	app, err := setup.SetupApp(config.Config{
		Port:      "7777",
		Dialector: &dl,
	}, SuccessEnv{})

	require.NoError(t, err)
	require.NotNil(t, app)
	require.NotNil(t, app.Server)
	require.NotNil(t, app.Services)
	require.Equal(t, "7777", app.Config.Port)
}
