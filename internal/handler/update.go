package handler

import (
	"net/http"
	"strings"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (h *Handler) Update(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain")

	if req.Method != http.MethodPost {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}
	if !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)

		return
	}

	path := strings.TrimPrefix(req.URL.Path, "/update/")
	values := strings.Split(path, "/")
	if len(values) != 3 {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	metric, err := models.NewMetricsFromStrings(values[0], values[1], values[2])
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	if err = h.Storage.Add(metric); err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	res.WriteHeader(http.StatusOK)
}
