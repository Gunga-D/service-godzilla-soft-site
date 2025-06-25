package subscription_product

import (
	"net/http"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/user"
)

type handler struct {
	userChecker userChecker
	subRepo     subRepo
}

func NewHandler(userChecker userChecker, subRepo subRepo) *handler {
	return &handler{
		userChecker: userChecker,
		subRepo:     subRepo,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID *int64
		if v, ok := r.Context().Value(user.MetaUserIDKey{}).(int64); ok {
			userID = pointer.ToInt64(v)
		}
		if userID == nil {
			api.Return401("Пользователь неавторизован", w)
			return
		}

		var body SubscriptionProductRequest
		if err := api.ReadBody(r, &body); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		hasSub, err := h.userChecker.HasSubscription(r.Context(), *userID)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if !hasSub {
			api.Return404("Подписка не подключена", w)
			return
		}

		product, err := h.subRepo.GetSubscriptionProduct(r.Context(), body.ItemID)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if product == nil {
			api.Return404("Аккаунта к данной игре не существует", w)
			return
		}

		api.ReturnOK(SubscriptionProductResponse{
			Login:    product.Login,
			Password: product.Password,
		}, w)
	}
}
