package services

import (
	"github.com/hotrungnhan/surl/models"
)

type GetShortenUrlParams struct {
	ID string
}

func (p *GetShortenUrlParams) ToFilter() *models.ShortenUrlFilterParams {
	return &models.ShortenUrlFilterParams{
		ID: &p.ID,
	}
}

type AddShortenUrlParams struct {
	OriginalUrl string
}

func (p *AddShortenUrlParams) ToCreateModel() *models.ShortenUrl {
	return &models.ShortenUrl{
		OriginalUrl: p.OriginalUrl,
	}
}

func (p *AddShortenUrlParams) ToFilter() *models.ShortenUrlFilterParams {
	return &models.ShortenUrlFilterParams{
		OriginalUrl: &p.OriginalUrl,
	}
}
