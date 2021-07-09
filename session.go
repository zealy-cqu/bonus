package bonus

import (
	"fmt"
	"net"
)

type session struct {
	*parser
	net.Conn
}

func newSession(d *dialer) (*session, error) {
	c, err := d.conn()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Success full connected.")
	p, err := newParser(c)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Success create parser")

	return &session{
		parser: p,
		Conn: c,
	}, nil
}
