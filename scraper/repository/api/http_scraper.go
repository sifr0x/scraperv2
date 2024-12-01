package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/datahuys/scraperv2/domain"
	"github.com/enriquebris/goconcurrentqueue"
)

type apiScraperRepository struct {
	queue *goconcurrentqueue.FIFO
}

func NewAPIScraperRepository() (domain.ScraperRepository, error) {
	repo := &apiScraperRepository{
		queue: goconcurrentqueue.NewFIFO(),
	}

	b, err := repo.fetchBatch()
	if err != nil {
		return nil, err
	}

	repo.addBatchToQueue(b)

	// go repo.enqueuer()

	return repo, nil
}

func (asr *apiScraperRepository) GetID() (int64, error) {
	item, err := asr.queue.DequeueOrWaitForNextElement()
	if err != nil {
		return 0, err
	}
	return item.(int64), nil
}

func (asr *apiScraperRepository) Close() error {
	return nil
}

func (asr *apiScraperRepository) enqueuer() {
	for {
		if asr.queue.GetLen() < 5000 {
			b, err := asr.fetchBatch()
			if err != nil {
				log.Printf("failed fetching batch: %v\n", err)
				continue
			}

			asr.addBatchToQueue(b)
		}

		time.Sleep(30 * time.Second)
	}
}

type batch struct {
	Content []int64 `json:"content"`
}

func (asr *apiScraperRepository) addBatchToQueue(b *batch) {
	for _, id := range b.Content {
		asr.queue.Enqueue(id)
	}

	log.Printf("added batch of %d items to queue\n", len(b.Content))
	log.Printf("queue length: %d\n", asr.queue.GetLen())
}

func (asr *apiScraperRepository) fetchBatch() (*batch, error) {
	res, err := http.Get("https://api-production-371e.up.railway.app/batch")
	if err != nil {
		return nil, err
	}

	// Decode JSON response into batch
	var b batch
	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		return nil, err
	}

	log.Printf("fetched batch of %d items\n", len(b.Content))

	return &b, nil
}
