package helpers

import (
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "" {
			name = fld.Tag.Get("xml")
		}
		if name == "" {
			name = fld.Tag.Get("form")
		}
		if name == "" {
			name = fld.Tag.Get("uri")
		}

		if name == "-" {
			return ""
		}

		return name
	})
}

func Validate(s any) (err error) {
	return validate.Struct(s)
}

func SetDefaults(s any) (err error) {
	err = defaults.Set(s)

	if err != nil {
		return err
	}
	return nil
}

func ValidateAndDefault(s any) (err error) {
	err = Validate(s)

	if err != nil {
		return err
	}

	err = SetDefaults(s)

	return nil
}
