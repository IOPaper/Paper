package core

import (
	"github.com/IOPaper/Paper/paper"
	"github.com/gin-gonic/gin"
)

func GetPaper(c *gin.Context) {
	paperIndex := c.Param("index")
	if paperIndex == "" {
		NewResponse(&Response{
			Code: 400,
			Msg:  "request param empty",
		}).Err(c)
		return
	}
	paperDoc, err := paper.Func.Find(paperIndex)
	if err != nil {
		NewResponse(&Response{
			Code: 400,
			Msg:  err.Error(),
		}).Err(c)
		return
	}
	NewResponse(&Response{
		Code: 200,
		Data: paperDoc.Export(),
	}).Ok(c)
}
