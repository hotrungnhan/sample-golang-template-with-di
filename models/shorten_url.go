package models

import (
	"github.com/hotrungnhan/surl/utils/types"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type ShortenUrlFilterParams struct {
	ID  *string   `filter:"eq;id"`
	IDs []*string `filter:"in;id"`

	OriginalUrl *string `filter:"eq;original_url"`
	// Preloads    []string
}

type ShortenUrl struct {
	types.BaseModel

	ID          string `gorm:"column:id;type:string;not null"`
	OriginalUrl string `gorm:"column:original_url;type:text;not null"`
}

// BeforeCreate GORM hook to set the ID
func (s *ShortenUrl) BeforeCreate(tx *gorm.DB) (err error) {
	const URL_SAFE_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	id, err := gonanoid.Generate(URL_SAFE_CHARS, 8)

	s.ID = id

	return err
}

func (*ShortenUrl) TableName() string {
	return "shorten_urls"
}
