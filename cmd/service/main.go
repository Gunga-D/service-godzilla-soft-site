package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	code_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/code/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/admin_create_item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/admin_warmup_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/cart_item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/categories_tree"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/create_order"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/item_details"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/mdw"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/payment_notification"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/user_login"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/user_register"
	item_cache "github.com/Gunga-D/service-godzilla-soft-site/internal/item/inmemory"
	item_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/item/postgres"
	order_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/order/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user/auth"
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

	databusClient := databus.NewClient(ctx)

	itemRepo := item_postgres.NewRepo(postgres)
	itemCache := item_cache.NewCache(itemRepo)

	userRepo := user_postgres.NewRepo(postgres)
	authJWT := auth.NewJwtService(os.Getenv("JWT_SECRET_KEY"))

	codeRepo := code_postgres.NewRepo(postgres)

	orderRepo := order_postgres.NewRepo(postgres)

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

		r1.Post("/user_register", user_register.NewHandler(authJWT, userRepo).Handle())
		r1.Post("/user_login", user_login.NewHandler(authJWT, userRepo).Handle())
		r1.Get("/item_details", item_details.NewHandler(itemCache).Handle())

		r1.Route("/", func(r2 chi.Router) {
			r2.Use(mdw.NewJWT(authJWT).VerifyUser)
			r2.Post("/cart_item", cart_item.NewHandler(codeRepo, itemCache, databusClient).Handle())
			r2.Post("/create_order", create_order.NewHandler(itemCache, orderRepo, databusClient).Handle())
		})

		r1.Post("/payment_notification", payment_notification.NewHandler().Handle())
	})
}
