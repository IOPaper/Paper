package impls

import (
	"github.com/IOPaper/Paper/paper/engine"
	"github.com/IOPaper/Paper/utils/result"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Paper interface {
	GetList(c *gin.Context)
	GetPaperIndexStatus(c *gin.Context)
	GetPaper(c *gin.Context)
	GetAttachment(c *gin.Context)
}

type paperImpl struct {
	engine.Engine
}

func NewImplPaper(e engine.Engine) Paper {
	return &paperImpl{Engine: e}
}

func (p *paperImpl) GetList(c *gin.Context) {
	var (
		b, l          = c.Query("before"), c.Query("limit")
		before, limit int
		err           error
	)
	if b != "" || l != "" {
		if before, err = strconv.Atoi(b); err != nil {
			result.New[any]().SetStatusCode(400).SetMessage("invalid request param").Err(c.Writer)
			return
		}
		if limit, err = strconv.Atoi(l); err != nil {
			result.New[any]().SetStatusCode(400).SetMessage("invalid request param").Err(c.Writer)
			return
		}
		if before < 0 || limit <= 0 {
			result.New[any]().SetStatusCode(400).SetMessage("invalid request param").Err(c.Writer)
			return
		}
	} else {
		before = 0
		limit = 10
	}
	list, err := p.GetPaperList(before, limit)
	if err != nil {
		result.New[gin.H]().SetStatusCode(503).SetMessage("can't fetch paper list").SetData(gin.H{
			"error": err.Error(),
		}).Err(c.Writer)
		return
	}
	result.New[*engine.PaperList]().
		SetStatusCode(200).
		SetMessage("ok").
		SetData(list).
		Ok(c.Writer)
}

func (p *paperImpl) GetPaperIndexStatus(c *gin.Context) {
	paperName := c.Param("index")
	if paperName == "" {
		result.New[any]().SetStatusCode(400).SetMessage("invalid request param").Err(c.Writer)
		return
	}
	result.New[gin.H]().
		SetStatusCode(200).
		SetMessage("ok").
		SetData(gin.H{
			"status": p.CheckPaperIndexStatus(paperName),
		}).
		Ok(c.Writer)
}

func (p *paperImpl) GetPaper(c *gin.Context) {
	paperName := c.Param("index")
	if paperName == "" {
		result.New[any]().SetStatusCode(400).SetMessage("invalid request param").Err(c.Writer)
		return
	}
	doc, err := p.GetOnePaper(paperName)
	if err != nil {
		result.New[gin.H]().SetStatusCode(503).SetMessage("can't get this paper").SetData(gin.H{
			"error": err.Error(),
		}).Err(c.Writer)
		return
	}
	result.New[*engine.Paper]().
		SetStatusCode(200).
		SetMessage("ok").
		SetData(doc).
		Ok(c.Writer)
}

func (p *paperImpl) GetAttachment(c *gin.Context) {
	paperName := c.Param("index")
	attachment := c.Param("attachment")
	if paperName == "" || attachment == "" || len(attachment) != 64 {
		result.New[any]().SetStatusCode(400).SetMessage("invalid request param").Err(c.Writer)
		return
	}
	doc, err := p.GetOnePaper(paperName)
	if err != nil {
		result.New[gin.H]().SetStatusCode(503).SetMessage("can't get this paper").SetData(gin.H{
			"error": err.Error(),
		}).Err(c.Writer)
		return
	}
	for _, index := range doc.Attachment {
		if index.Hash == attachment {
			path, er := p.GetAttachmentPath(attachment)
			if er != nil {
				result.New[gin.H]().SetStatusCode(503).SetMessage("can't get this attachment").SetData(gin.H{
					"error": err.Error(),
				})
				return
			}
			c.File(path)
			return
		}
	}
	result.New[any]().SetStatusCode(404).SetMessage("attachment not found").Err(c.Writer)
}
