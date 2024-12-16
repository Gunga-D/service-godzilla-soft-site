package admin_save_thumbnail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
)

const (
	selfDiskURL = "https://disk.godzillasoft.ru/upload/%s"
)

type handler struct {
	httpClient   *http.Client
	diskLogin    string
	diskPassword string
}

func NewHandler(diskLogin string, diskPassword string) *handler {
	return &handler{
		httpClient:   &http.Client{},
		diskLogin:    diskLogin,
		diskPassword: diskPassword,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AdminSaveThumbnailRequest
		if err := api.ReadBody(r, &req); err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}

		data, err := base64.StdEncoding.DecodeString(req.DataBase64)
		if err != nil {
			api.Return400(err.Error(), w)
			return
		}
		buff := &bytes.Buffer{}
		buff.WriteString(string(data))

		diskReq, err := http.NewRequest("PUT", fmt.Sprintf(selfDiskURL, req.FileName), buff)
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}

		diskReq.SetBasicAuth(h.diskLogin, h.diskPassword)
		diskResp, err := h.httpClient.Do(diskReq)
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}
		diskRespBody, err := io.ReadAll(diskResp.Body)
		if err != nil {
			api.Return500(err.Error(), w)
			return
		}
		w.Write(diskRespBody)
	}
}
