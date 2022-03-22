package consts

const (
	GET_CONTESTS_BY_TIME_INTERVAL ApiCall = iota
)

var (
	domain = "http://localhost:8080/"
)

type ApiCall int

func (a ApiCall) Addr() string {
	return a.AddDomain(a.Path())
}

func (a ApiCall) Path() string {
	return []string{"contests/by-time-interval"}[a]
}

func (a ApiCall) AddDomain(path string) string {
	return domain + path
}
