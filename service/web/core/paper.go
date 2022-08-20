package core

import (
	"github.com/IOPaper/Paper/paper"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetPaperList(c *gin.Context) {
	var (
		before, limit int
		err           error
	)
	if before, err = strconv.Atoi(c.Query("before")); err != nil {
		NewResponse(&Response{Code: 400, Msg: "invalid request param"}).Err(c)
		return
	}
	if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
		NewResponse(&Response{Code: 400, Msg: "invalid request param"}).Err(c)
		return
	}
	if before < 0 || limit <= 0 || limit > 10 {
		NewResponse(&Response{Code: 400, Msg: "invalid request param"}).Err(c)
		return
	}

	list, err := paper.Func.List(uint(before), uint(limit))
	if err != nil {
		NewResponse(&Response{Code: 400, Msg: err.Error()}).Err(c)
		return
	}

	NewResponse(&Response{
		Code: 200,
		Data: list.Export(),
	}).Ok(c)

}

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

func GetPaperAttachment(c *gin.Context) {

}
