package domain

type UserAgentRepository interface {
    GetRandomUserAgent() (agent string, err error)
}
