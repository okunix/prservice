package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/okunix/prservice/internal/app/config"
	"github.com/okunix/prservice/internal/pkg/data"
	"github.com/okunix/prservice/internal/pkg/migrations"
)

func Run() error {
	conf := config.GetConfig()
	err := data.InitPostgres(conf.PostgresConfig, migrations.MigrationsFS)
	if err != nil {
		return err
	}

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("healthy"))
	})

	server := http.Server{
		Addr:    conf.Addr,
		Handler: router,
	}
	return server.ListenAndServe()
}
