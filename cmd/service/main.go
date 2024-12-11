package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/admin_create_item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/admin_warmup_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/categories_tree"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/mdw"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/user_login"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/user_register"
	item_cache "github.com/Gunga-D/service-godzilla-soft-site/internal/item/inmemory"
	item_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/item/postgres"
	user_password "github.com/Gunga-D/service-godzilla-soft-site/internal/user/password"
	user_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/user/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	postgres := postgres.Get(ctx)

	itemRepo := item_postgres.NewRepo(postgres)
	itemCache := item_cache.NewCache(itemRepo)

	userRepo := user_postgres.NewRepo(postgres)

	userPwdGenerator := user_password.NewGenerator()

	mux := chi.NewMux()
	mux.Use(middleware.Timeout(5 * time.Second))
	c := cors.New(cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowedHeaders:      []string{"*"},
		AllowPrivateNetwork: true,
		AllowCredentials:    true,
		Debug:               false,
	})
	mux.Use(c.Handler)

	mux.Route("/api/v1", func(r1 chi.Router) {
		r1.Route("/admin", func(r2 chi.Router) {
			r2.Post("/admin_warmup_items", admin_warmup_items.NewHandler(itemCache).Handle())
			r2.Post("/create_item", admin_create_item.NewHandler(itemRepo).Handle())
		})

		r1.Get("/categories_tree", categories_tree.NewHandler().Handle())

		r1.Post("/user_register", user_register.NewHandler(os.Getenv("JWT_SECRET_KEY"), userRepo, userPwdGenerator).Handle())
		r1.Post("/user_login", user_login.NewHandler(os.Getenv("JWT_SECRET_KEY"), userRepo, userPwdGenerator).Handle())

		r1.Route("/", func(r2 chi.Router) {
			r2.Use(mdw.NewJWT(os.Getenv("JWT_SECRET_KEY")).VerifyUser)
		})
	})
}
