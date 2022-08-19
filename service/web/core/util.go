package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

func structFieldEmptyCheck(v any) error {
	typeof := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	for i := 0; i < typeof.NumField(); i++ {
		if typeof.Field(i).Tag.Get("check") != "" {
			if value.Field(i).IsZero() {
				return fmt.Errorf("field %s is empty", typeof.Field(i).Name)
			}
		}
	}
	return nil
}

type SysCreateNewPaperRequest struct {
	Title       string   `json:"title" check:"title"`
	ContentSize uint     `json:"content_size" check:"content_size"`
	Tags        []string `json:"tags,omitempty"`
	// AttachmentCount uint     `json:"attachment_count,omitempty"`
	Author string `json:"author" check:"author"`
}
type SysPutNewPaperRequest struct {
	Id      int64  `json:"paper_id" check:"paper_id"`
	Content string `json:"content" check:"content"`
}

func (r *SysCreateNewPaperRequest) EmptyField() error {
	return structFieldEmptyCheck(*r)
}

func (r *SysPutNewPaperRequest) EmptyField() error {
	return structFieldEmptyCheck(*r)
}

////////////////////////////////////////////////////////////////////////////

type h = map[string]any

type Response struct {
	Code int
	Msg  string
	Data any
}

type Request struct {
	ctx *gin.Context
}

type Result struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg,omitempty"`
	Data   any    `json:"data,omitempty"`
}

type ResponseResult interface {
	Ok(c *gin.Context)
	Err(c *gin.Context)
}

type RequestResult interface {
	JSON(v any) error
}

func NewRequest(c *gin.Context) RequestResult {
	return &Request{ctx: c}
}

func NewResponse(r *Response) ResponseResult {
	return r
}

func (r *Response) Ok(c *gin.Context) {
	c.JSON(r.Code, Result{
		Status: true,
		Msg:    r.Msg,
		Data:   r.Data,
	})
}

func (r *Response) Err(c *gin.Context) {
	c.JSON(r.Code, Result{
		Status: false,
		Msg:    r.Msg,
		Data:   r.Data,
	})
}

func (r *Request) JSON(v any) error {
	if err := r.ctx.BindJSON(v); err != nil {
		NewResponse(&Response{
			Code: 400,
			Msg:  fmt.Sprintf("parse request data fail, %s", err.Error()),
		}).Err(r.ctx)
		return err
	}
	return nil
}
