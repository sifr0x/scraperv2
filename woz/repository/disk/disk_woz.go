package disk

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/datahuys/scraperv2/domain"
)

type diskWozRepository struct {
	main string
}

func NewDiskWozRepository(main string) domain.WozStorageRepository {
	repo := &diskWozRepository{main}
	_ = os.Mkdir(repo.main, os.ModePerm)
	return repo
}

func (d *diskWozRepository) GetByID(id int64) (wozList []domain.Woz, err error) {
	base := fmt.Sprintf("%s/%d", d.main, id)
	files, err := os.ReadDir(base)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		split := strings.Split(file.Name(), ".")
		i, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			return wozList, err
		}
		tm := time.Unix(i, 0)

		body, err := os.ReadFile(base + "/" + file.Name())
		if err != nil {
			return wozList, err
		}

		woz := domain.Woz{
			ID:        id,
			Payload:   body,
			ScrapedAt: tm,
		}
		wozList = append(wozList, woz)
	}

	return
}

func (d *diskWozRepository) Store(woz domain.Woz) (err error) {
	base := fmt.Sprintf("%s/%d", d.main, woz.ID)
	_ = os.Mkdir(base, os.ModePerm)

	path := fmt.Sprintf("%s/%d.json", base, woz.ScrapedAt.Unix())
	err = os.WriteFile(path, woz.Payload, 0644)
	return
}
