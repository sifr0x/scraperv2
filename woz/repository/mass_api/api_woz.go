package massapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/datahuys/scraperv2/domain"
)

type apiWozRepository struct {
	endpoint string
	client   *http.Client
}

func NewApiWozRepository(endpoint string, client *http.Client) domain.WozStorageRepository {
	return &apiWozRepository{endpoint, client}
}

func (a *apiWozRepository) GetByID(id int64) (wozList []domain.Woz, err error) {
	err = errors.New("not implemented")
	return
}

func (a *apiWozRepository) Store(woz domain.Woz) (err error) {
	if woz.Status == 404 {
		err = a.StoreMissing(woz)
		return
	}

	var wrapper Data
	err = json.Unmarshal(woz.Payload, &wrapper)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal payload: %w", err)
		return
	}

	if wrapper.WOZObject.NummerAanduidingID == 0 {
		wrapper.WOZObject.NummerAanduidingID = woz.ID
		woz.Payload, err = json.Marshal(wrapper)
		if err != nil {
			err = fmt.Errorf("failed to re-marshal payload: %w", err)
			return
		}
	}

	req, err := http.NewRequest("POST", a.endpoint, bytes.NewBuffer(woz.Payload))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		bodyBts, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyBts))
		fmt.Println(string(woz.Payload))
		fmt.Println(woz)

		err = fmt.Errorf("unexpected status code: %d", res.StatusCode)
		return
	}

	return
}

func (a *apiWozRepository) StoreMissing(woz domain.Woz) (err error) {
	missing := domain.MissingObject{
		BagID:     woz.ID,
		CheckedAt: woz.ScrapedAt,
	}

	marshal, err := json.Marshal(missing)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", a.endpoint+"/missing", bytes.NewBuffer(marshal))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		bodyBts, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyBts))
		fmt.Println(string(marshal))

		err = fmt.Errorf("unexpected missing status code: %d", res.StatusCode)
		return
	}
	return
}
