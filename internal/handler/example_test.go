package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/handler"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	"github.com/vrnvgasu/metrics/internal/service/audit"
	"github.com/vrnvgasu/metrics/internal/service/metric"
)

func newExampleRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	publisher, _ := audit.NewAudit(&config.ServerCnf{})
	repo := mem.NewMemStorage()
	s := metric.NewService(repo)
	h := handler.NewHandler(s, nil, publisher, &config.ServerCnf{})
	return handler.NewRouter(h)
}

// ExampleHandler_Update демонстрирует обновление метрики через URL-параметры.
// POST /update/:mtype/:name/:value
func ExampleHandler_Update() {
	router := newExampleRouter()

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/123.5", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output: 200
}

// ExampleHandler_UpdateJSON демонстрирует обновление одной метрики через JSON.
// POST /update/
func ExampleHandler_UpdateJSON() {
	router := newExampleRouter()

	body := `{"id":"PollCount","type":"counter","delta":5}`
	req := httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output: 200
}

// ExampleHandler_UpdateJSONList демонстрирует пакетное обновление метрик.
// POST /updates
func ExampleHandler_UpdateJSONList() {
	router := newExampleRouter()

	body := `[{"id":"Alloc","type":"gauge","value":1.5},{"id":"PollCount","type":"counter","delta":1}]`
	req := httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output: 200
}

// ExampleHandler_Find демонстрирует получение значения метрики текстом.
// GET /value/:mtype/:name
func ExampleHandler_Find() {
	router := newExampleRouter()

	// сначала сохраняем метрику
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/123.5", http.NoBody)
	router.ServeHTTP(httptest.NewRecorder(), req)

	// получаем значение
	req = httptest.NewRequest(http.MethodGet, "/value/gauge/Alloc", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())
	// Output:
	// 200
	// 123.5
}

// ExampleHandler_Value демонстрирует получение метрики через JSON.
// POST /value/
func ExampleHandler_Value() {
	router := newExampleRouter()

	// сначала сохраняем метрику
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/1.5", http.NoBody)
	router.ServeHTTP(httptest.NewRecorder(), req)

	// получаем значение через JSON
	body := `{"id":"Alloc","type":"gauge"}`
	req = httptest.NewRequest(http.MethodPost, "/value/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output: 200
}

// ExampleHandler_List демонстрирует получение списка всех метрик в HTML.
// GET /
func ExampleHandler_List() {
	router := newExampleRouter()

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output: 200
}
