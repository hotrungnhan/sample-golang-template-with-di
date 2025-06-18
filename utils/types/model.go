package types

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	// ID        ulid.ULID `gorm:"primaryKey;type:ulid;default:gen_ulid()" json:"id"`
	CreatedAt time.Time
}

type UpdatableModel struct {
	UpdatedAt time.Time
}
type SoftDeletableModel struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
