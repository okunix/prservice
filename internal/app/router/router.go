package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/okunix/prservice/internal/app/endpoints"
	"github.com/okunix/prservice/internal/pkg/data"
	"github.com/okunix/prservice/internal/pkg/repos"
	"github.com/okunix/prservice/static"
)

func New() http.Handler {
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

	r := chi.NewRouter()
	r.Use(corsMiddleware.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.FileServerFS(static.StaticFS))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("healthy"))
	})

	db := data.PostgreSQL()
	userRepo := repos.NewUserRepo(db)
	teamRepo := repos.NewTeamRepo(db, userRepo)
	//pullRequestRepo := repos.NewPullRequestRepo(db)

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", endpoints.AddTeam(teamRepo))
		r.Get("/get", nil)
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/getReview", nil)
		r.Post("/setIsActive", nil)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Get("/create", nil)
		r.Get("/merge", nil)
		r.Get("/reassign", nil)
	})

	return r
}
