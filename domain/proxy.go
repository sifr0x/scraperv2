package domain

type ProxyRepository interface {
	GetRandomUserProxy() (proxy string, err error)
}
