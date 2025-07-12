package controllers

import (
	"github.com/hotrungnhan/go-automapper"
	"github.com/hotrungnhan/surl/services"
)

type GetShortenUrlParams struct {
	ID string `uri:"id" validate:"required,len=8"`
}

type AddShortenUrlParams struct {
	OriginalUrl string `json:"original_url" xml:"original_url" form:"original_url" validate:"required,url"`
}

func init() {
	mapper.Register(mapper.Global, func(src GetShortenUrlParams) services.GetShortenUrlParams {
		return services.GetShortenUrlParams{
			ID: src.ID,
		}
	})
	mapper.Register(mapper.Global, func(src AddShortenUrlParams) services.AddShortenUrlParams {
		return services.AddShortenUrlParams{
			OriginalUrl: src.OriginalUrl,
		}
	})
}
