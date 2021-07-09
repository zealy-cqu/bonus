package bonus

import (
	"crypto/rc4"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
)

const pkgLen = 2

type parser struct {
	dec *decoder
	enc *encoder
	conn net.Conn
}

func newParser(conn net.Conn) (*parser, error) {
	fmt.Printf("Before create parser.\n")
	ps := &parser{
		conn: conn,
	}

	fmt.Printf("remote addr: %v\n", conn.RemoteAddr())
	_, msg, err := ps.read(false)
	if err !=nil {
		return nil, err
	}

	fmt.Printf("Succesfully parser.read.\n")
	keys := msg.(*CipherKeyNtf)
	fmt.Printf("CipherKeyNtf:%+#v\n", keys)
	decipher, err := rc4.NewCipher([]byte(keys.SvrKey))
	if err != nil {
		return nil, err
	}
	ps.dec = newDecoder(decipher)

	encipher, err := rc4.NewCipher([]byte(keys.CltKey))
	if err != nil {
		return nil, err
	}
	ps.enc = newEncoder(encipher)

	return ps, nil
}


func (p *parser) read(crypt bool) (CSID, proto.Message, error)  {
	fmt.Printf("Before start parser.read.\n")
	lenb := make([]byte, pkgLen)
	if _, err := io.ReadFull(p.conn, lenb); err != nil {
		return 0, nil, err
	}

	lenn := binary.BigEndian.Uint16(lenb)

	payload := make([]byte, lenn)
	if _, err := io.ReadFull(p.conn, payload); err != nil {
		return 0, nil, err
	}

	fmt.Printf("")
	return p.dec.decode(payload, crypt)
}

func (p *parser) write(msgID CSID, msg proto.Message, crypt bool) error {
	payload, err := p.enc.encode(msgID, msg, crypt)
	if err != nil {
		return nil
	}

	pkg := make([]byte, pkgLen)
	binary.BigEndian.PutUint16(pkg, uint16(len(payload)))

	pkg = append(pkg, payload...)

	if _, err := p.conn.Write(pkg); err != nil {
		return err
	}
	return nil
}

