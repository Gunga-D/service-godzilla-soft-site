package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/template"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/tinkoff"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/yandex_market"
	code_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/code/postgres"
	collection_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/collection/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/add_review"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/admin_create_item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/admin_load_codes"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/admin_recalc_price"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/admin_save_thumbnail"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/admin_warmup_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/cart_item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/categories_tree"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/check_voucher"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/collection_details"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/create_order"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/fetch_collection_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/fetch_collections"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/fetch_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/item_details"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/mdw"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/new_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/payment_notification"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/popular_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/recomendation_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/reviews"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/sales_items"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/search_suggest"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/sitemap"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/steam_calc_price"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/steam_gift_resolve_profile"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/steam_invoice"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/telegram_sign_in"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/think"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/think_result"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/user_change_password"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/user_login"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/user_profile"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/http/user_register"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item/fillers"
	steam_filler "github.com/Gunga-D/service-godzilla-soft-site/internal/item/fillers/steam"
	yandex_market_filler "github.com/Gunga-D/service-godzilla-soft-site/internal/item/fillers/yandex_market"
	item_cache "github.com/Gunga-D/service-godzilla-soft-site/internal/item/inmemory"
	item_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/item/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item/recomendation"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item/suggest"
	neuro_redis "github.com/Gunga-D/service-godzilla-soft-site/internal/neuro/redis"
	order_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/order/postgres"
	review_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/review/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user/auth"
	user_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/user/postgres"
	voucher_activation "github.com/Gunga-D/service-godzilla-soft-site/internal/voucher/activation"
	voucher_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/voucher/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/redis"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/service"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/transport/listener"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	tele "gopkg.in/telebot.v4"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	postgres := postgres.Get(ctx)
	redis := redis.Get(ctx)

	databusClient := databus.NewClient(ctx)

	telebot, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Printf("[error] cant init telegram bot: %v", err)
		return
	}
	logger.Get().SetBot(telebot)

	steamClient := steam.NewClient(os.Getenv("STEAM_KEY"))
	tinkoffClient := tinkoff.NewClient(os.Getenv("TINKOFF_API_URL"), os.Getenv("TINKOFF_PASSWORD"), os.Getenv("TINKOFF_TERMINAL_KEY"))
	yaMarketBusinessID, err := strconv.Atoi(os.Getenv("YA_MARKET_BUSINESS_ID"))
	if err != nil {
		log.Printf("[error] cant get yander marker business id: %v", err)
		return
	}
	yaMarket := yandex_market.NewClient(os.Getenv("YA_MARKET_API_URL"), os.Getenv("YA_MARKET_AUTH"), int64(yaMarketBusinessID))

	changePasswordTemplate, err := template.ParseFiles("assets/change-password-template.html")
	if err != nil {
		log.Fatalln("failed loading html templates")
	}

	itemRepo := item_postgres.NewRepo(postgres)
	itemRecommendation := recomendation.NewService(steamClient)
	itemCache := item_cache.NewCache(itemRepo, []fillers.Filler{
		yandex_market_filler.NewFiller(yaMarket),
		steam_filler.NewFiller(),
	}, itemRecommendation)
	neuroCache := neuro_redis.NewRepo(redis)
	go itemCache.StartSync(ctx)
	itemSuggestSrv := suggest.NewService(itemRepo, itemCache)
	go itemSuggestSrv.StartSync(ctx)

	userRepo := user_postgres.NewRepo(postgres, redis)
	authJWT := auth.NewJwtService(os.Getenv("JWT_SECRET_KEY"))

	voucherRepo := voucher_postgres.NewRepo(postgres)
	voucherActivation := voucher_activation.NewService(voucherRepo)

	codeRepo := code_postgres.NewRepo(postgres)
	orderRepo := order_postgres.NewRepo(postgres)
	collectionRepo := collection_postgres.NewRepo(postgres)

	reviewRepo := review_postgres.NewRepo(postgres)

	mux := chi.NewMux()
	mux.Use(middleware.Timeout(5 * time.Second))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD",
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		Debug: false,
	})
	mux.Use(c.Handler)

	mux.Route("/v1", func(r1 chi.Router) {
		r1.Use(mdw.NewUseragent().IdentifyPlatform)

		r1.Get("/sitemap.xml", sitemap.NewHandler(itemCache).Handle())

		r1.Route("/admin", func(r2 chi.Router) {
			r2.Use(mdw.NewBearerMDW(os.Getenv("ADMIN_SECRET_KEY")).VerifyUser)
			r2.Post("/warmup_items", admin_warmup_items.NewHandler(itemCache).Handle())
			r2.Post("/create_item", admin_create_item.NewHandler(itemRepo).Handle())
			r2.Post("/load_codes", admin_load_codes.NewHandler(codeRepo, itemRepo, databusClient).Handle())
			r2.Post("/save_thumbnail", admin_save_thumbnail.NewHandler(os.Getenv("GODZILLA_SOFT_DISK_LOGIN"), os.Getenv("GODZILLA_SOFT_DISK_PASSWORD")).Handle())
			r2.Post("/recalc_price", admin_recalc_price.NewHandler(itemRepo, steamClient).Handle())
		})

		r1.Get("/categories_tree", categories_tree.NewHandler().Handle())

		r1.Post("/search_suggest", search_suggest.NewHandler(itemSuggestSrv).Handle())
		r1.Post("/user_register", user_register.NewHandler(authJWT, userRepo).Handle())
		r1.Post("/user_login", user_login.NewHandler(authJWT, userRepo).Handle())
		r1.Post("/user_change_password", user_change_password.NewHandler(userRepo, changePasswordTemplate, databusClient).Handle())
		r1.Post("/telegram_sign_in", telegram_sign_in.NewHandler(authJWT, userRepo, os.Getenv("TELEGRAM_LOGIN_WIDGET_BOT_TOKEN")).Handle())

		r1.Get("/collections", fetch_collections.NewHandler(collectionRepo).Handle())
		r1.Get("/collection_items", fetch_collection_items.NewHandler(itemCache, collectionRepo).Handle())
		r1.Get("/collection_details", collection_details.NewHandler(collectionRepo).Handle())

		r1.Get("/popular_items", popular_items.NewHandler(itemCache).Handle())
		r1.Get("/recomendation_items", recomendation_items.NewHandler(itemCache).Handle())
		r1.Get("/sales_items", sales_items.NewHandler(itemCache).Handle())
		r1.Get("/new_items", new_items.NewHandler(itemCache).Handle())
		r1.Get("/items", fetch_items.NewHandler(itemRepo, itemCache).Handle())
		r1.Get("/item_details", item_details.NewHandler(itemCache, itemRecommendation).Handle())

		// Пополнение Steam
		r1.Post("/steam_calc", steam_calc_price.NewHandler().Handle())
		r1.Post("/steam_invoice", steam_invoice.NewHandler(orderRepo, tinkoffClient).Handle())

		r1.Get("/reviews", reviews.NewHandler(reviewRepo).Handle())

		r1.Route("/", func(r2 chi.Router) {
			r2.Use(mdw.NewJWT(authJWT).VerifyUser)
			r2.Post("/check_voucher", check_voucher.NewHandler(itemCache, voucherActivation).Handle())
			r2.Post("/cart_item", cart_item.NewHandler(codeRepo, itemCache, databusClient).Handle())
			r2.Post("/create_order", create_order.NewHandler(itemCache, orderRepo, tinkoffClient, databusClient, voucherActivation).Handle())

			r2.Post("/add_review", add_review.NewHandler(reviewRepo).Handle())

			r2.Get("/user_profile", user_profile.NewHandler(redis, userRepo).Handle())
		})

		r1.Route("/steam_gift", func(r2 chi.Router) {
			r2.Post("/resolve_profile", steam_gift_resolve_profile.NewHandler(steamClient).Handle())
		})

		r1.Post("/payment_notification", payment_notification.NewHandler(os.Getenv("TINKOFF_TERMINAL_KEY"), orderRepo).Handle())

		r1.Post("/think", think.NewHandler(databusClient).Handle())
		r1.Post("/think_result", think_result.NewHandler(neuroCache).Handle())
	})

	log.Println("[info] server start up")
	err = service.Listen(ctx, listener.NewHTTP(), mux)
	if err != nil {
		log.Printf("[error] server finished with an error: %v", err)
		return
	}
}
