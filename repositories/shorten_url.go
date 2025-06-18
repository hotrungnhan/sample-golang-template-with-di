package repositories

import (
	"context"
	"errors"

	"github.com/hotrungnhan/surl/generated/queries"
	"github.com/hotrungnhan/surl/models"
	"github.com/hotrungnhan/surl/utils/injects"
	"github.com/samber/do/v2"

	"gorm.io/gorm"
)

type ShortenUrlRepository interface {
	Get(context.Context, *models.ShortenUrlFilterParams) (*models.ShortenUrl, error)
	Add(context.Context, *models.ShortenUrl) (*models.ShortenUrl, error)
}

type shortenUrlRepositoryImpl struct {
	cache  *injects.Cache
	config injects.ApplicationConfig
	BaseRepository
}

func (r *shortenUrlRepositoryImpl) Get(ctx context.Context, params *models.ShortenUrlFilterParams) (*models.ShortenUrl, error) {
	db := r.GetReadDB(ctx)

	db, err := r.BuildQuery(db, params)

	if err != nil {
		return nil, err
	}
	record, err := queries.Use(db).ShortenUrl.First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return record, nil
}

func (r *shortenUrlRepositoryImpl) Add(ctx context.Context, record *models.ShortenUrl) (*models.ShortenUrl, error) {
	err := queries.Use(r.GetWriteDB(ctx)).ShortenUrl.Create(record)
	return record, err
}

func NewShortenUrlRepository(i do.Injector) (ShortenUrlRepository, error) {
	return &shortenUrlRepositoryImpl{
		BaseRepository: newBaseRepository(do.MustInvoke[*injects.DB](i)),
	}, nil
}
