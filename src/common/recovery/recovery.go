package recovery

import (
	"context"
	"github.com/go-errors/errors"
	"kits/api/src/common/fault"
	"kits/api/src/common/logger"
)

func HandleRoutine() {
	if err := recover(); err != nil {
		goErr := errors.Wrap(err, 2)
		wErr := fault.Wrapf(goErr, "[Recovery] go routine")
		logger.Fatal(context.Background(), wErr, "[Recovery] stack %s", goErr.Stack())
	}
}
