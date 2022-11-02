package impls

import (
	"github.com/IOPaper/Paper/crypto"
	"github.com/IOPaper/Paper/paper/engine"
	"github.com/gin-gonic/gin"
)

type Writing interface {
	AddPaper(c *gin.Context)
}

type writingImpl struct {
	secret string
	*crypto.EcdsaKeypair
	engine.Engine
}

func NewImplWriting(e engine.Engine, keypair *crypto.EcdsaKeypair, secret string) Writing {
	return &writingImpl{Engine: e, EcdsaKeypair: keypair, secret: secret}
}

func (p *writingImpl) AddPaper(c *gin.Context) {

}
