package bonus

import (
	"github.com/golang/protobuf/proto"
)

type message struct {
	id  CSID
	msg proto.Message
}

type writeMsg struct {
	msg   *message
	errCh chan error
}

type readMsg struct {
	msg *message
	err error
}

type writeCh chan *writeMsg
type readCh chan *readMsg

type closeCh chan struct{}
