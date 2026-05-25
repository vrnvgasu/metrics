// Package helper - пакет со вспомогательными для тестирования функциями.
package helper

// NewRefFloat64 получить указатель на float64.
func NewRefFloat64(v float64) *float64 {
	return &v
}

// NewRefInt64 получить указатель на int64.
func NewRefInt64(v int64) *int64 {
	return &v
}
