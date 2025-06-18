package helpers

import (
	"errors"

	"github.com/hotrungnhan/surl/utils/types"

	"github.com/gofiber/fiber/v3"
)

// Bind Validate and Set Defaults for parameters in a Fiber context

type BindType uint8

const (
	Body   = 1 << iota                            // 00000001
	Cookie                                        // 00000010
	Header                                        // 00000100
	Query                                         // 00001000
	URI                                           // 00010000
	All    = Body | Cookie | Header | Query | URI // 00011111
)

func BindValidateDefaultCtx(ctx fiber.Ctx, BindType BindType, params any) error {

	if BindType&Cookie != 0 {
		if err := ctx.Bind().Cookie(params); err != nil {
			return types.BadRequestError.
				WithCode("DATA_PARSER").
				WithError(errors.Join(errors.New("Bind Cookie Error"), err))
		}
	}

	if BindType&Header != 0 {
		if err := ctx.Bind().Cookie(params); err != nil {
			return types.BadRequestError.
				WithCode("DATA_PARSER").
				WithError(errors.Join(errors.New("Bind Header Error"), err))
		}
	}

	if BindType&Query != 0 {
		if err := ctx.Bind().Query(params); err != nil {
			return types.BadRequestError.
				WithCode("DATA_PARSER").
				WithError(errors.Join(errors.New("Bind Query Error"), err))
		}
	}

	if BindType&Body != 0 && len(ctx.BodyRaw()) > 0 {
		if err := ctx.Bind().Body(params); err != nil {

			return types.BadRequestError.
				WithCode("DATA_PARSER").
				WithError(errors.Join(errors.New("Bind Body Error"), err))
		}
	}
	if BindType&URI != 0 {
		if err := ctx.Bind().URI(params); err != nil {
			return types.BadRequestError.
				WithCode("DATA_PARSER").
				WithError(errors.Join(errors.New("Bind URI Error"), err))
		}
	}
	if err := Validate(params); err != nil {
		return types.BadRequestError.
			WithError(err)
	}

	if err := SetDefaults(params); err != nil {
		return types.InternalServerError.
			WithError(err)
	}

	return nil
}
