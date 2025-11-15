package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/okunix/prservice/internal/app/config"
	"github.com/okunix/prservice/internal/app/router"
	"github.com/okunix/prservice/internal/pkg/data"
	"github.com/okunix/prservice/migrations"
)

func Run() error {
	conf := config.GetConfig()
	err := data.InitPostgres(conf.PostgresConfig, migrations.MigrationsFS)
	if err != nil {
		return err
	}

	if conf.AdminToken == config.DefaultAdminToken {
		fmt.Fprintf(os.Stderr, "WARNING: Default admin token is used!")
	}

	server := http.Server{
		Addr:    conf.Addr,
		Handler: router.New(),
	}
	return server.ListenAndServe()
}
