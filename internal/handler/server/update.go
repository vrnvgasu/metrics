package server

import (
	"net/http"
	"strings"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func Update(res http.ResponseWriter, req *http.Request) {
	//if req.Method != http.MethodPost {
	//	http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	//
	//	return
	//}
	//if !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
	//	http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
	//
	//	return
	//}

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

	if err = repository.Storage.Add(metric); err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	res.Header().Add("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
}
