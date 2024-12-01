package main

import (
	"log"

	"github.com/datahuys/scraperv2/proxy/repository/mullvad"
	"github.com/datahuys/scraperv2/scraper/repository/file"
	"github.com/datahuys/scraperv2/scraper/usecase"
	"github.com/datahuys/scraperv2/user_agent/repository/local"
	"github.com/datahuys/scraperv2/woz/repository/disk"
	"github.com/datahuys/scraperv2/woz/repository/loket"
)

func main() {
	scraperRepo, err := file.NewFileScraperRepository("ids.txt", "pos.txt")
	// scraperRepo, err := _scraperRepo.NewAPIScraperRepository()
	if err != nil {
		log.Fatalln(err)
	}
	defer scraperRepo.Close()

	diskRepo := disk.NewDiskWozRepository("woz_data")

	repo, err := local.NewLocalUserAgentRepository("agents.json")
	if err != nil {
		log.Fatalln(err)
	}

	mvRepo, err := mullvad.NewMullvadProxyRepository("https://api-www.mullvad.net/www/relays/all/")
	if err != nil {
		log.Fatalln(err)
	}

	scraperUcase := usecase.NewScraperUsecase(scraperRepo, mvRepo, repo, diskRepo, loket.NewLoketWozRepository)
	scraperUcase.Start()

}
