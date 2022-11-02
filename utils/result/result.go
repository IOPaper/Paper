package result

import (
	"github.com/json-iterator/go"
	"io"
	"net/http"
)

type Result[T any] struct {
	statusCode int
	Status     bool   `json:"status"`
	Message    string `json:"msg,omitempty"`
	Data       T      `json:"data,omitempty"`
}

type ResponseEncoding interface {
	JSON() []byte
	JSONPipe(w io.Writer)
}

type Setter[T any] interface {
	SetStatusCode(code int) Setter[T]
	SetMessage(msg string) Setter[T]
	SetData(data T) Setter[T]
	Encode() ResponseEncoding
	Err(w http.ResponseWriter)
	Ok(w http.ResponseWriter)
}

type Encoding[T any] struct {
	result *Result[T]
}

func (r *Encoding[T]) JSON() []byte {
	result, err := jsoniter.ConfigFastest.Marshal(r.result)
	if err != nil {
		return nil
	}
	return result
}

func (r *Encoding[T]) JSONPipe(w io.Writer) {
	jsoniter.ConfigFastest.NewEncoder(w).Encode(r.result)
}

func NewResult[T any]() Setter[T] {
	return &Result[T]{statusCode: http.StatusOK}
}

func (r *Result[T]) SetStatusCode(code int) Setter[T] {
	r.statusCode = code
	return r
}

func (r *Result[T]) SetMessage(msg string) Setter[T] {
	r.Message = msg
	return r
}

func (r *Result[T]) SetData(data T) Setter[T] {
	r.Data = data
	return r
}

func (r *Result[T]) GetStatus() bool {
	return r.Status
}

func (r *Result[T]) GetMessage() string {
	return r.Message
}

func (r *Result[T]) GetData() T {
	return r.Data
}

func (r *Result[T]) Encode() ResponseEncoding {
	return &Encoding[T]{
		result: r,
	}
}

func (r *Result[T]) Err(w http.ResponseWriter) {
	w.WriteHeader(r.statusCode)
	w.Header().Set("Content-Type", "application/json")
	r.Status = false
	r.Encode().JSONPipe(w)
}

func (r *Result[T]) Ok(w http.ResponseWriter) {
	w.WriteHeader(r.statusCode)
	w.Header().Set("Content-Type", "application/json")
	r.Status = true
	r.Encode().JSONPipe(w)
}
