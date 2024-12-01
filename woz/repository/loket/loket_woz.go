package loket

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/datahuys/scraperv2/domain"
	"golang.org/x/net/publicsuffix"
)

type loketWozRepository struct {
	client *http.Client
	agent  string
}

func NewLoketWozRepository(ctx context.Context, client *http.Client, agent string) (domain.WozRepository, error) {
	var err error
	repo := &loketWozRepository{client, agent}
	err = repo.startSession(ctx)
	return repo, err
}

func (lw *loketWozRepository) startSession(ctx context.Context) (err error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	lw.client.Jar = jar

	_, _, err = lw.execute(ctx, "POST", "session/start")
	return
}

func (lw *loketWozRepository) GetByID(ctx context.Context, id int64) (woz domain.Woz, err error) {
	action := fmt.Sprintf("wozwaarde/nummeraanduiding/%016d", id)
	res, status, err := lw.execute(ctx, "GET", action)
	if err != nil {
		return
	}
	woz.ID = id
	woz.Payload = res
	woz.Status = status
	woz.ScrapedAt = time.Now()
	return
}

func (lw *loketWozRepository) execute(ctx context.Context, method string, action string) (res []byte, status int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // Ensure resources are cleaned up

	url := fmt.Sprintf("https://wozwaardeloket.nl/wozwaardeloket-api/v1/%s", action)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", lw.agent)
	resp, err := lw.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	res, err = ioutil.ReadAll(resp.Body)
	status = resp.StatusCode
	return
}

func (lw *loketWozRepository) Store(woz domain.Woz) (err error) {
	err = errors.New("not implemented")
	return
}
