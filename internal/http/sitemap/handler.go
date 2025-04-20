package sitemap

import (
	"fmt"
	"net/http"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
)

type handler struct {
	getter getter
}

func NewHandler(getter getter) *handler {
	return &handler{
		getter: getter,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm := stm.NewSitemap(0)

		// Дефолтные категории
		sm.Create()
		sm.Add(stm.URL{{"loc", "https://godzillasoft.ru/steam_deposit"}, {"changefreq", "daily"}, {"priority", "0.8"}})
		sm.Add(stm.URL{{"loc", "https://godzillasoft.ru/random"}, {"changefreq", "daily"}, {"priority", "0.8"}})
		sm.Add(stm.URL{{"loc", "https://godzillasoft.ru/games"}, {"changefreq", "daily"}, {"priority", "0.8"}})
		sm.Add(stm.URL{{"loc", "https://godzillasoft.ru/deposits"}, {"changefreq", "daily"}, {"priority", "0.8"}})
		sm.Add(stm.URL{{"loc", "https://godzillasoft.ru/subscriptions"}, {"changefreq", "daily"}, {"priority", "0.8"}})

		// Блог
		sm.Add(stm.URL{{"loc", "https://godzillasoft.ru/blog/luchshie-igri-2024-goda"}, {"changefreq", "daily"}, {"priority", "0.8"}})
		sm.Add(stm.URL{{"loc", "https://godzillasoft.ru/blog/zelenie-magazin-ili-stoit-li-pokupat-lizenzirovanie-igri"}, {"changefreq", "daily"}, {"priority", "0.8"}})

		items, err := h.getter.FetchAllItems(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, i := range items {
			var catalogPath string
			switch i.CategoryID {
			case 10001:
				catalogPath = "games"
			case 10002:
				catalogPath = "subscriptions"
			case 10004:
				catalogPath = "deposits"
			default:
				continue
			}
			sm.Add(stm.URL{{"loc", fmt.Sprintf("https://godzillasoft.ru/%s/%s", catalogPath, generatePathValue(i.Title, i.ID))}, {"changefreq", "daily"}, {"priority", "0.8"}})
		}

		w.Write(sm.XMLContent())
	}
}
