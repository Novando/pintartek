package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func Validate(data interface{}) error {
	validate := validator.New()
	errs := validate.Struct(data)
	if errs != nil {
		var errMsgs []string
		for _, err := range errs.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.Field(),
				err.Value(),
				err.Tag(),
			))
		}
		return fmt.Errorf(strings.Join(errMsgs, "; "))
	}
	return nil
}
