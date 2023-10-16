package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"kits/api/src/common/fault"
	"kits/api/src/common/logger"
	httputil "kits/api/src/present/http/util"
)

func Recovery(c *gin.Context, err any) {
	goErr := errors.Wrap(err, 2)
	wErr := fault.Wrapf(goErr, "[Recovery] go gin")

	logger.Fatal(c.Request.Context(), wErr, "[Recovery] stack %s", goErr.Stack())
	httputil.ServeErrResponse(c, wErr)
	c.Abort()
}
