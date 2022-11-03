package impls

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/IOPaper/Paper/crypto"
	"github.com/IOPaper/Paper/paper/engine"
	"github.com/IOPaper/Paper/utils/result"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
)

type Writing interface {
	GetPaperIndexStatus(c *gin.Context)
	AddPaper(c *gin.Context)
}

type writingImpl struct {
	secret            string
	maxAttachmentSize int64
	*crypto.EcdsaKeypair
	engine.Engine
}

func NewImplWriting(e engine.Engine, keypair *crypto.EcdsaKeypair, secret string, maxAttachmentSize int64) Writing {
	return &writingImpl{Engine: e, EcdsaKeypair: keypair, secret: secret}
}

func (p *writingImpl) GetPaperIndexStatus(c *gin.Context) {
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

func (p *writingImpl) AddPaper(c *gin.Context) {
	var dto engine.PaperDTO
	if err := c.BindJSON(&dto); err != nil {
		result.New[any]().SetStatusCode(400).SetMessage("invalid request").Err(c.Writer)
		return
	}
	verify := false
	if dto.Sign != nil {
		verify = p.Verify(bytes.NewBufferString(dto.Body), dto.Sign)
	}
	if err := p.AddOnePaper(verify, &dto); err != nil {
		result.New[any]().SetStatusCode(200).SetMessage(err.Error()).Err(c.Writer)
		return
	}
	result.New[gin.H]().
		SetStatusCode(201).
		SetMessage("ok").
		SetData(gin.H{
			"verify": verify,
		}).
		Ok(c.Writer)
}

func (p *writingImpl) UploadPaperAttachment(c *gin.Context) {

	paperName := c.Param("index")
	attachment := c.Param("attachment")
	if paperName == "" || attachment == "" || len(attachment) != 64 {
		result.New[any]().SetStatusCode(400).SetMessage("invalid request param").Err(c.Writer)
		return
	}

	file, err := c.FormFile("attachment")
	switch {
	case err != nil:
		result.New[any]().
			SetStatusCode(400).
			SetMessage(errors.Wrap(err, "attachment upload fail, error:").Error()).
			Err(c.Writer)
		return
	case file.Size > p.maxAttachmentSize:
		result.New[any]().SetStatusCode(503).SetMessage("attachment size too large").Err(c.Writer)
		return
	}

	f, err := file.Open()
	if err != nil {
		result.New[any]().SetStatusCode(400).SetMessage(err.Error()).Err(c.Writer)
		return
	}

	// copy file to byte slice
	defer f.Close()
	var b []byte
	if _, err = io.ReadFull(f, b); err != nil {
		result.New[any]().SetStatusCode(503).SetMessage("copy file error").Err(c.Writer)
		return
	}

	// sum file sha256 and equal
	{
		s := sha256.New()
		io.Copy(s, bytes.NewBuffer(b))
		if attachment != hex.EncodeToString(s.Sum(nil)) {
			result.New[any]().SetStatusCode(400).SetMessage("invalid attachment").Err(c.Writer)
			return
		}
	}

	if err = p.AddAttachment(attachment, bytes.NewBuffer(b)); err != nil {
		result.New[any]().
			SetStatusCode(503).
			SetMessage(errors.Wrap(err, "attachment upload fail, error:").Error()).
			Err(c.Writer)
		return
	}

	// todo: set to nil to trigger GC quickly
	b = nil

	result.New[any]().SetStatusCode(201).SetMessage("ok").Ok(c.Writer)

}
