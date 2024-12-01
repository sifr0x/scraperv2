package local

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/datahuys/scraperv2/domain"
	"github.com/fsnotify/fsnotify"
)

type localUserAgentRepository struct {
	path   string
	agents atomic.Value
}

func NewLocalUserAgentRepository(path string) (domain.UserAgentRepository, error) {
	repo := &localUserAgentRepository{path: path}
	agents, err := repo.fetch()
	if err != nil {
		return repo, err
	}

	repo.agents.Store(agents)

	go repo.watcher()

	return repo, nil
}

func (l *localUserAgentRepository) fetch() (agents []string, err error) {
	agentsFile, err := os.ReadFile(l.path)
	if err != nil {
		return
	}

	err = json.Unmarshal(agentsFile, &agents)
	log.Printf("Fetched %d agents\n", len(agents))
	return
}

func (l *localUserAgentRepository) watcher() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
					agents, err := l.fetch()
					if err != nil {
						log.Println("failed fetching modified file:", err)
						return
					}

					l.agents.Swap(agents)
					log.Println("swapped modified file")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(l.path)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func (l *localUserAgentRepository) GetRandomUserAgent() (agent string, err error) {
	agentsAtomic := l.agents.Load()
	if agentsAtomic == nil {
		err = errors.New("no agents available")
		return
	}

	agents := agentsAtomic.([]string)

	len := len(agents)
	n := uint32(0)
	if len > 0 {
		n = getRandomUint32() % uint32(len)
	}
	agent = agents[n]
	return
}

func getRandomUint32() uint32 {
	x := time.Now().UnixNano()
	return uint32((x >> 32) ^ x)
}
