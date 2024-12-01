package mullvad

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/datahuys/scraperv2/domain"
)

type mullvadProxyRepository struct {
	path    string
	proxies []string
}

type Server struct {
	Hostname             string        `json:"hostname"`
	CountryCode          string        `json:"country_code"`
	CountryName          string        `json:"country_name"`
	CityCode             string        `json:"city_code"`
	CityName             string        `json:"city_name"`
	Active               bool          `json:"active"`
	Owned                bool          `json:"owned"`
	Provider             string        `json:"provider"`
	Ipv4AddrIn           string        `json:"ipv4_addr_in"`
	Ipv6AddrIn           string        `json:"ipv6_addr_in"`
	NetworkPortSpeed     int           `json:"network_port_speed"`
	Stboot               bool          `json:"stboot"`
	Type                 string        `json:"type"`
	StatusMessages       []interface{} `json:"status_messages"`
	Pubkey               string        `json:"pubkey,omitempty"`
	MultihopPort         int           `json:"multihop_port,omitempty"`
	SocksName            string        `json:"socks_name,omitempty"`
	SocksPort            int           `json:"socks_port,omitempty"`
	Ipv4V2Ray            string        `json:"ipv4_v2ray,omitempty"`
	SSHFingerprintSha256 string        `json:"ssh_fingerprint_sha256,omitempty"`
	SSHFingerprintMd5    string        `json:"ssh_fingerprint_md5,omitempty"`
}

func NewMullvadProxyRepository(path string) (domain.ProxyRepository, error) {
	var err error
	repo := &mullvadProxyRepository{path: path}
	repo.proxies, err = repo.fetch()
	log.Printf("Fetched %d proxies\n", len(repo.proxies))
	return repo, err
}

func (m *mullvadProxyRepository) fetch() (proxies []string, err error) {
	resp, err := http.Get(m.path)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var servers []Server
	err = json.Unmarshal(bts, &servers)
	if err != nil {
		return
	}

	for _, server := range servers {
		if server.SocksName == "" || server.SocksPort == 0 {
			continue
		}
		proxy := fmt.Sprintf("%s:%d", server.SocksName, server.SocksPort)
		proxies = append(proxies, proxy)
	}
	return
}

func (m *mullvadProxyRepository) GetRandomUserProxy() (proxy string, err error) {
	lent := len(m.proxies)
	if lent == 0 {
		err = errors.New("no proxies available")
		return
	}

	n := uint32(0)
	if lent > 0 {
		n = getRandomUint32() % uint32(lent)
	}
	proxy = m.proxies[n]
	return
}

func getRandomUint32() uint32 {
	x := time.Now().UnixNano()
	return uint32((x >> 32) ^ x)
}
