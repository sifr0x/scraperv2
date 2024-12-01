package domain

type ScraperUsecase interface {
	Start()
}

type ScraperRepository interface {
	GetID() (int64, error)
	Close() error
}
