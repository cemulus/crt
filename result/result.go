package result

type QueryResult interface {
	Table() []byte
	JSON() ([]byte, error)
	CSV() ([]byte, error)
	Size() int
}
