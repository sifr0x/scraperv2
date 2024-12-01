package usecase

import (
	"context"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/datahuys/scraperv2/domain"
	"golang.org/x/net/proxy"

	_massRepo "github.com/datahuys/scraperv2/woz/repository/mass_api"

	"golang.org/x/sync/semaphore"
)

type scraperUsecase struct {
	scraperRepo domain.ScraperRepository
	proxyRepo   domain.ProxyRepository
	agentRepo   domain.UserAgentRepository
	wozDiskRepo domain.WozStorageRepository
	wozFunc     func(context.Context, *http.Client, string) (domain.WozRepository, error)

	wg  *sync.WaitGroup
	sem *semaphore.Weighted
}

func NewScraperUsecase(scraperRepo domain.ScraperRepository, proxyRepo domain.ProxyRepository, agentRepo domain.UserAgentRepository, wozDiskRepo domain.WozStorageRepository, wozFunc func(context.Context, *http.Client, string) (domain.WozRepository, error)) domain.ScraperUsecase {
	wg := new(sync.WaitGroup)
	sem := semaphore.NewWeighted(48)

	return &scraperUsecase{
		scraperRepo,
		proxyRepo,
		agentRepo,
		wozDiskRepo,
		wozFunc,
		wg,
		sem,
	}
}

func (u *scraperUsecase) worker(num int) {
	defer u.sem.Release(1)

	err := u.newSession()
	if err != nil {
		log.Println(num, "new session worker error:", err)
		log.Println(u.sem)
	}
}

func (u *scraperUsecase) Start() {
	i := 0
	ctx := context.TODO()
	for {
		if err := u.sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		go u.worker(i)
		i++
	}
}

func (u *scraperUsecase) newSession() (err error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	// defer cancel() // Ensure resources are cleaned up
	ctx := context.Background()

	scrapedItems := 0

	agent, err := u.agentRepo.GetRandomUserAgent()
	if err != nil {
		return
	}

	proxyStr, err := u.proxyRepo.GetRandomUserProxy()
	if err != nil {
		return
	}

	baseDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	dialSocksProxy, err := proxy.SOCKS5("tcp", proxyStr, nil, baseDialer)
	if err != nil {
		return
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial:                  dialSocksProxy.Dial,
			MaxIdleConns:          10,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}

	// massRepo := _massRepo.NewApiWozRepository("https://internal.huisapi.nl/woz", httpClient)
	massRepo := _massRepo.NewApiWozRepository("https://api-production-371e.up.railway.app/woz", &http.Client{})

	session, err := u.wozFunc(ctx, httpClient, agent)
	if err != nil {
		return
	}

	for scrapedItems != 10 {
		id, err := u.scraperRepo.GetID()
		if err != nil {
			return err
		}

		woz, err := session.GetByID(ctx, id)
		if err != nil {
			return err
		}

		err = massRepo.Store(woz)
		if err != nil {
			log.Printf("failed to store woz: %s\n", err)
			err = u.wozDiskRepo.Store(woz)
			if err != nil {
				return err
			}
		}

		scrapedItems++
	}
	log.Printf("scraped %d items with %s/%s\n", scrapedItems, proxyStr, agent)
	return
}
