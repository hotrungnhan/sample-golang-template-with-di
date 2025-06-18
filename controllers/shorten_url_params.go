package controllers

import (
	"github.com/hotrungnhan/surl/services"
)

type GetShortenUrlParams struct {
	ID string `uri:"id" validate:"required,len=8"`
}

func (p *GetShortenUrlParams) ToServiceParams() *services.GetShortenUrlParams {
	return &services.GetShortenUrlParams{
		ID: p.ID,
	}
}

type AddShortenUrlParams struct {
	OriginalUrl string `json:"original_url" xml:"original_url" form:"original_url" validate:"required,url"`
}

func (p *AddShortenUrlParams) ToServiceParams() *services.AddShortenUrlParams {
	return &services.AddShortenUrlParams{
		OriginalUrl: p.OriginalUrl,
	}
}
