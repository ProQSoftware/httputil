package httputil

import (
	"io"
	"net/http"
)

// HttpError - интерфейс абстрагирующий поведение всех
// ошибок получаемых получаемых в результате совершения
// http запросов и их обработке
type HttpError interface {
	error
	Status() int
	ContentType() string
}

// Error - базовая ошибка
type Error struct {
	msg         string
	statusCode  int
	contentType string
}

// WriteHttpError функция записывающая ошибку в ответ
func WriteHttpError(w http.ResponseWriter, e error) {
	switch e.(type) {
	case HttpError:
		err := e.(HttpError)

		w.Header().Add("Content-Type", err.ContentType())
		w.WriteHeader(err.Status())
		_, _ = io.WriteString(w, err.Error())
		// http.Error(w, err.Error(), err.Status())
	default:
		http.Error(w, e.Error(), http.StatusTeapot)
	}
}

// Error - метод реализующий интерфейс error
// Нужен для того чтобы достать описание ошибки
func (e *Error) Error() string {
	return e.msg
}

// Status - метод реализующий интерфейс HttpError
// Нужен чтобы указать верный статус при записи ошибки
func (e *Error) Status() int {
	return e.statusCode
}

// ContentType - метод реализующий интерфейс HttpError
// Нужен для того чтобы указать при записи ошибки в ответ
// тип содержимого этой ошибка (json или html и т д)
func (e *Error) ContentType() string {
	return e.contentType
}

// NewError принимает первым аргументом статус выполнения, вторым аргументом сообщение об ошибке
func NewError(status int, contentType, msg string) *Error {
	return &Error{
		msg:         msg,
		statusCode:  status,
		contentType: contentType,
	}
}
