package serializers

import (
	"time"

	mapper "github.com/hotrungnhan/go-automapper"
	"github.com/hotrungnhan/surl/models"
)

type ShortUrlSerializer struct {
	ID          string `json:"id"`
	OriginalUrl string `json:"original_url"`
}

type ShortUrlDetailSerializer struct {
	ID          string `json:"id"`
	OriginalUrl string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
}

func init() {
	mapper.Register(mapper.Global, func(record models.ShortenUrl) ShortUrlDetailSerializer {
		return ShortUrlDetailSerializer{
			ID:          record.ID,
			OriginalUrl: record.OriginalUrl,
			CreatedAt:   record.CreatedAt.Format(time.RFC3339),
		}
	})
	mapper.Register(mapper.Global, func(record models.ShortenUrl) ShortUrlSerializer {
		return ShortUrlSerializer{
			ID:          record.ID,
			OriginalUrl: record.OriginalUrl,
		}
	})
}
