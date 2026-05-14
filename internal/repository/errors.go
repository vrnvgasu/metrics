package repository

import "errors"

var (
	ErrNotFound   = errors.New("not found")   // метрика не найдена
	ErrNotSupport = errors.New("not support") // операция не поддерживается хранилищем
)
