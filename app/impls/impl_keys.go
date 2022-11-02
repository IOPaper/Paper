package impls

import (
	"github.com/IOPaper/Paper/crypto"
	"github.com/IOPaper/Paper/utils/result"
	"github.com/gin-gonic/gin"
)

type Keys interface {
	Public(c *gin.Context)
}

type keys struct {
	*crypto.EcdsaKeypair
}

func NewImplKeys(keypair *crypto.EcdsaKeypair) Keys {
	return &keys{EcdsaKeypair: keypair}
}

func (k *keys) Public(c *gin.Context) {
	_, pub := k.ExportKeypair()
	result.New[gin.H]().
		SetStatusCode(200).
		SetMessage("ok").
		SetData(gin.H{
			"public_key": pub,
		}).
		Ok(c.Writer)
}
