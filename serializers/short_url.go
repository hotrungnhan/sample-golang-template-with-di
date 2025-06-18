package serializers

import (
	"time"

	"github.com/hotrungnhan/surl/models"
	"github.com/hotrungnhan/surl/utils/types"

	"github.com/samber/lo"
)

type ShortUrlSerializer struct {
	ID          string `json:"id"`
	OriginalUrl string `json:"original_url"`
}

func NewShortUrlSerializer(record *models.ShortenUrl) types.ISerializer {
	return &ShortUrlSerializer{
		ID:          record.ID,
		OriginalUrl: record.OriginalUrl,
	}
}

func NewListShortUrlSerializer(records []*models.ShortenUrl) types.ISerializer {
	return lo.Map(records, func(record *models.ShortenUrl, _ int) types.ISerializer {
		return NewShortUrlSerializer(record)
	})
}

type ShortUrlDetailSerializer struct {
	ID          string `json:"id"`
	OriginalUrl string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
}

func NewShortUrlDetailSerializer(record *models.ShortenUrl) types.ISerializer {
	return &ShortUrlDetailSerializer{
		ID:          record.ID,
		OriginalUrl: record.OriginalUrl,
		CreatedAt:   record.CreatedAt.Format(time.RFC3339),
	}
}

func NewListShortUrlDetailSerializer(records []*models.ShortenUrl) types.ISerializer {
	return lo.Map(records, func(record *models.ShortenUrl, _ int) types.ISerializer {
		return NewShortUrlDetailSerializer(record)
	})
}
