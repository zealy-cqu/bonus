package bonus

import (
	"crypto/rc4"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
)

const msgIndexLen = 2

var deEntry = decodeEntry{
	CSID_ID_CipherKeyNtf: decodeCipherKeyNtf,
	CSID_ID_HeartbeatRsp: decodeHeartbeatRsp,
}

type decoder struct {
	decipher *rc4.Cipher
	//deEntry decodeEntry
}

type decodeFn func(int []byte) (proto.Message, error)
type decodeEntry map[CSID]decodeFn

func newDecoder(decipher *rc4.Cipher) *decoder {
	return &decoder{
		decipher: decipher,
		//deEntry: deEntry,
	}
}

func (d *decoder) decode(in []byte, crypt bool) (CSID, proto.Message, error) {
	payload := in
	if crypt {
		payload = make([]byte, len(in))
		d.decipher.XORKeyStream(payload, in)
	}

	msgID := binary.BigEndian.Uint16(payload[:msgIndexLen])
	fmt.Printf("decoder: %+#v", d)
	msg, err := deEntry[CSID(msgID)](payload[msgIndexLen:])
	if err != nil {
		return CSID(msgID), nil, err
	}
	return CSID(msgID), msg, nil
}

func decodeCipherKeyNtf(in []byte) (proto.Message, error) {
	msg := &CipherKeyNtf{}
	if err := proto.Unmarshal(in, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func decodeHeartbeatRsp(in []byte) (proto.Message, error) {
	msg := &HeartbeatRsp{}
	if err := proto.Unmarshal(in, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
