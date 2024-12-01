package domain

import "time"

type MissingObject struct {
	BagID     int64     `db:"bag_id" json:"bag_id"`
	CheckedAt time.Time `db:"checked_at" json:"checked_at"`
}

type MissingObjectRepository interface {
	Store(missingObject MissingObject) error
}
