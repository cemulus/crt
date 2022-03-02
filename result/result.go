package result

type QueryResult interface {
	Table() string
	JSON() (string, error)
	CSV() (string, error)
	Size() int
}
