package domain

import (
	"context"
	"time"
)

type Woz struct {
	ID        int64
	Payload   []byte
	Status    int
	ScrapedAt time.Time
}

type WozRepository interface {
	GetByID(ctx context.Context, id int64) (Woz, error)
	Store(Woz) error
}

type WozStorageRepository interface {
	GetByID(id int64) ([]Woz, error)
	Store(Woz) error
}
