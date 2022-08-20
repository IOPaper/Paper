package core

import (
	"github.com/IOPaper/Paper/paper"
	"github.com/gin-gonic/gin"
)

func CreateNewPaper(c *gin.Context) {
	req := &SysCreateNewPaperRequest{}
	err := NewRequest(c).JSON(&req)
	if err != nil {
		return
	}
	paperIndex := c.Param("index")
	if paperIndex == "" {
		NewResponse(&Response{Code: 401, Msg: "invalid request"}).Err(c)
		return
	}
	if err = req.EmptyField(); err != nil {
		NewResponse(&Response{Code: 400, Msg: err.Error()}).Err(c)
		return
	}
	create, err := paper.Func.Create(paperIndex)
	if err != nil {
		NewResponse(&Response{Code: 400, Msg: err.Error()}).Err(c)
		return
	}
	// TODO: put create request metadata
	create.SetContentSize(req.ContentSize)
	create.SetTitle(req.Title)
	create.SetTags(req.Tags...)
	create.SetAuthor(req.Author)
	NewResponse(&Response{
		Code: 201,
		Data: h{
			"paper_path": paperIndex,
			"paper_id":   create.GetPaperId(),
		},
	}).Ok(c)
	return
}

func PutNewPaper(c *gin.Context) {
	req := &SysPutNewPaperRequest{}
	err := NewRequest(c).JSON(&req)
	if err != nil {
		return
	}
	paperIndex := c.Param("index")
	if paperIndex == "" {
		NewResponse(&Response{Code: 401, Msg: "invalid request"}).Err(c)
		return
	}
	if err = req.EmptyField(); err != nil {
		NewResponse(&Response{Code: 400, Msg: err.Error()}).Err(c)
		return
	}
	create, err := paper.Func.RecoverCreate(req.Id, paperIndex)
	if err != nil {
		NewResponse(&Response{Code: 400, Msg: err.Error()}).Err(c)
		return
	}
	if err = create.SetContent(req.Content); err != nil {
		NewResponse(&Response{Code: 400, Msg: err.Error()}).Err(c)
		return
	}
	if err = create.Done(); err != nil {
		NewResponse(&Response{Code: 400, Msg: err.Error()}).Err(c)
		return
	}
	NewResponse(&Response{Code: 201}).Ok(c)
}

func PutAttachment(c *gin.Context) {

}
