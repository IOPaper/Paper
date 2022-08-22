package module

import "github.com/gin-gonic/gin"

type Call interface {
	Action() gin.HandlerFunc
}

type Calls []Call

func newCalls() Calls {
	return make(Calls, 0)
}

var (
	mds Calls = newCalls()
)

func registerCall(call Call) {
	mds = append(mds, call)
}

func Setup(e *gin.Engine) {
	if len(mds) == 0 {
		return
	}
	for i := 0; i < len(mds); i++ {
		e.Use(mds[i].Action())
	}
}
