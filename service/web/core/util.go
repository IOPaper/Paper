package core

import "github.com/gin-gonic/gin"

type h = map[string]any

type Response struct {
	Code int
	Msg  string
	Data any
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
