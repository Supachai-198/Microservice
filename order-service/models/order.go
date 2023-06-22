package models

import (
	"time"

	"github.com/lib/pq"
)

type Order struct {
	ID        uint          `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	UserId    uint          `json:"user_id" gorm:"type:integer;not null"`
	ProductId pq.Int64Array `json:"product_id" gorm:"type:integer[];not null"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"update_at"`
}
