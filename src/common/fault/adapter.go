package fault

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	appvalidator "kits/api/src/common/validator"
	"strings"
)

func DBWrapf(err error, format string, v ...interface{}) *AppError {
	msg := fmt.Sprintf(format, v...)
	wErr := fmt.Errorf("%s: %w", msg, err)
	tag := getTagFromErr(err)

	return &AppError{tag: tag, cause: wErr}
}

func getTagFromErr(err error) Tag {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return TagNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return TagAlreadyExists
	}
	if errors.Is(err, context.Canceled) {
		return TagCancelled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return TagDeadlineExceeded
	}

	return TagInternal
}

func ConvertValidatorErr(err error) *AppError {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)

	if !ok {
		return Wrap(err).SetTag(TagInvalidArgument)
	}

	trans := appvalidator.GetTranslator()

	messages := make([]string, 0, len(errs))

	for _, e := range errs {
		messages = append(messages, e.Translate(trans))
	}
	msg := strings.Join(messages, "; ")

	return Wrap(fmt.Errorf(msg)).SetTag(TagInvalidArgument).SetMessage(msg)
}
