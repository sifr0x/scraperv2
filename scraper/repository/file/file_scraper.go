package file

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/datahuys/scraperv2/domain"
)

type fileScraperRepository struct {
	filePath string
	posPath  string
	ids      []int64

	pos          *int64
	lastSavedPos int64
}

func NewFileScraperRepository(filePath string, posPath string) (domain.ScraperRepository, error) {
	var err error
	repo := &fileScraperRepository{
		filePath: filePath,
		posPath:  posPath,

		lastSavedPos: -1,
	}
	posNew := int64(-1)
	repo.pos = &posNew

	repo.ids, err = repo.fetch()
	if err != nil {
		return repo, err
	}
	log.Printf("loaded %d ids\n", len(repo.ids))

	pos, err := repo.fetchPos()
	if err != nil {
		log.Printf("failed loading pos, using default value: %v\n", err)
		repo.lastSavedPos = -1
	} else {
		log.Printf("loaded pos at %d\n", pos)
		repo.lastSavedPos = pos
	}
	atomic.StoreInt64(repo.pos, repo.lastSavedPos)

	go repo.syncRoutine()

	return repo, nil
}

func (fr *fileScraperRepository) fetch() (ids []int64, err error) {
	f, err := os.Open(fr.filePath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		scanned := scanner.Text()
		id, err := strconv.ParseInt(scanned, 10, 64)
		if err != nil {
			log.Printf("not id %s: %v", scanned, err)
			continue
		}
		ids = append(ids, id)
	}

	err = scanner.Err()
	return
}

func (fr *fileScraperRepository) fetchPos() (pos int64, err error) {
	data, err := os.ReadFile(fr.posPath)
	if err != nil {
		return
	}

	posStr := string(data)
	posStr = strings.TrimSpace(posStr)

	pos, err = strconv.ParseInt(posStr, 10, 64)
	return
}

func (fr *fileScraperRepository) syncPos() (err error) {
	current := atomic.LoadInt64(fr.pos)
	if fr.lastSavedPos == current {
		return
	}

	posStr := fmt.Sprintf("%d", current)
	data := []byte(posStr)
	err = os.WriteFile(fr.posPath, data, 0644)
	if err == nil {
		fr.lastSavedPos = current
	} else {
		log.Println(err)
	}
	return
}

func (fr *fileScraperRepository) Close() (err error) {
	return fr.syncPos()
}

func (fr *fileScraperRepository) syncRoutine() {
	for {
		time.Sleep(5 * time.Second)
		err := fr.syncPos()
		if err != nil {
			log.Printf("failed syncing pos: %v\n", err)
		}
	}
}

func (f *fileScraperRepository) GetID() (id int64, err error) {
	pos := atomic.AddInt64(f.pos, 1)

	if pos >= int64(len(f.ids)) {
		err = errors.New("no ids available")
		return
	}

	id = f.ids[pos]
	return
}
