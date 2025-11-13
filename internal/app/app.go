package app

import (
	"net/http"

	"github.com/okunix/prservice/internal/app/config"
	"github.com/okunix/prservice/internal/app/router"
	"github.com/okunix/prservice/internal/pkg/data"
	"github.com/okunix/prservice/internal/pkg/migrations"
)

func Run() error {
	conf := config.GetConfig()
	err := data.InitPostgres(conf.PostgresConfig, migrations.MigrationsFS)
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:    conf.Addr,
		Handler: router.New(),
	}
	return server.ListenAndServe()
}
