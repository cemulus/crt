package result

type Printer interface {
	Table() []byte
	JSON() ([]byte, error)
	CSV() ([]byte, error)
	Size() int
}
