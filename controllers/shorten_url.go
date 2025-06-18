package controllers

import (
	"net/http"

	"github.com/hotrungnhan/surl/serializers"
	"github.com/hotrungnhan/surl/services"
	"github.com/hotrungnhan/surl/utils/helpers"
	"github.com/hotrungnhan/surl/utils/types"

	"github.com/gofiber/fiber/v3"
	"github.com/samber/do/v2"
)

type ShortenUrlController interface {
	types.Controller

	Get(ctx fiber.Ctx) error
	Add(ctx fiber.Ctx) error
}

type shortenUrlControllerImpl struct {
	shortenUrlService services.ShortenUrlService
	di                do.Injector
}

func NewShortenUrlController(di do.Injector) (ShortenUrlController, error) {
	return &shortenUrlControllerImpl{
		di:                di,
		shortenUrlService: do.MustInvoke[services.ShortenUrlService](di),
	}, nil
}
func (c *shortenUrlControllerImpl) GetPublic(ctx fiber.Ctx) error {
	params := &GetShortenUrlParams{}

	if err := helpers.BindValidateDefaultCtx(ctx, helpers.Query|helpers.URI, params); err != nil {
		return err
	}

	data, err := c.shortenUrlService.Get(ctx.Context(), params.ToServiceParams())

	if err != nil {
		return err
	}

	ctx.Set(fiber.HeaderLocation, data.(*serializers.ShortUrlDetailSerializer).OriginalUrl)
	return ctx.SendStatus(http.StatusFound)
}

func (c *shortenUrlControllerImpl) Get(ctx fiber.Ctx) error {
	params := &GetShortenUrlParams{}

	if err := helpers.BindValidateDefaultCtx(ctx, helpers.Query|helpers.URI, params); err != nil {
		return err
	}

	data, err := c.shortenUrlService.Get(ctx.Context(), params.ToServiceParams())

	if err != nil {
		return err
	}

	return types.OKResponse.WithData(data)
}

func (c *shortenUrlControllerImpl) Add(ctx fiber.Ctx) error {
	params := &AddShortenUrlParams{}

	if err := helpers.BindValidateDefaultCtx(ctx, helpers.Body, params); err != nil {
		return err
	}
	data, err := c.shortenUrlService.Add(ctx.Context(), params.ToServiceParams())

	if err != nil {
		return err
	}

	return types.AcceptedResponse.WithData(data).WithStatusCode(http.StatusCreated)
}

func (c *shortenUrlControllerImpl) RegisterEndpoint(router fiber.Router) {

	router.Get("/shortlinks/:id", c.GetPublic).Name("Redirect Shorten Url")

	apiShortLink := router.Group("/api/shortlinks")

	apiShortLink.Get("/:id", c.Get).Name("Get Shorten Url")
	apiShortLink.Post("/", c.Add).Name("Add Shorten Url")

}
