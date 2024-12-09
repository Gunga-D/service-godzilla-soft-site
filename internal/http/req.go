package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func ReadBody(req *http.Request, m interface{}) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	return json.Unmarshal(body, m)
}
