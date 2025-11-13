package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/okunix/prservice/internal/app/config"
	"github.com/okunix/prservice/internal/pkg/data"
	"github.com/okunix/prservice/internal/pkg/migrations"
	"github.com/okunix/prservice/static"
)

func Run() error {
	conf := config.GetConfig()
	err := data.InitPostgres(conf.PostgresConfig, migrations.MigrationsFS)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{
			"*",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router := chi.NewRouter()
	router.Use(corsMiddleware.Handler)
	router.Handle("/static/*", http.FileServerFS(static.StaticFS))

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
