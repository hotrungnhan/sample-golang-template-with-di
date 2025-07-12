package services

import (
	"context"

	mapper "github.com/hotrungnhan/go-automapper"
	"github.com/hotrungnhan/surl/models"
	"github.com/hotrungnhan/surl/repositories"
	"github.com/hotrungnhan/surl/serializers"
	"github.com/hotrungnhan/surl/utils/injects"
	"github.com/hotrungnhan/surl/utils/types"

	"github.com/samber/do/v2"
)

type ShortenUrlService interface {
	Get(ctx context.Context, params *GetShortenUrlParams) (types.ISerializer, error)
	Add(ctx context.Context, params *AddShortenUrlParams) (types.ISerializer, error)
}
type shortenUrlServiceImpl struct {
	shortenRepo repositories.ShortenUrlRepository
	config      injects.ApplicationConfig
}

func (s *shortenUrlServiceImpl) Get(ctx context.Context, params *GetShortenUrlParams) (types.ISerializer, error) {

	record, err := s.shortenRepo.Get(ctx, params.ToFilter())
	if err != nil {
		return nil, types.InternalServerError
	}

	if record == nil {
		return nil, types.NoContentError
	}

	return mapper.MustMap[*models.ShortenUrl, *serializers.ShortUrlDetailSerializer](mapper.Global, record), nil
}

func (s *shortenUrlServiceImpl) Add(ctx context.Context, params *AddShortenUrlParams) (types.ISerializer, error) {
	record, err := s.shortenRepo.Get(ctx, params.ToFilter())

	if err != nil {
		return nil, types.InternalServerError
	}

	if record == nil {
		record, err = s.shortenRepo.Add(ctx, params.ToCreateModel())
		if err != nil {
			return nil, types.InternalServerError.WithErrorString("Error when create shorten url")
		}
	}

	return mapper.MustMap[*models.ShortenUrl, *serializers.ShortUrlSerializer](mapper.Global, record), nil
}

func NewShortenUrlService(di do.Injector) (ShortenUrlService, error) {
	return &shortenUrlServiceImpl{
		shortenRepo: do.MustInvoke[repositories.ShortenUrlRepository](di),
	}, nil
}
