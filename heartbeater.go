package bonus

import (
	"sync"
)

// heartbeater is responsible for reading and writing protobuf messages
type heartbeater struct {
	wg      *sync.WaitGroup
	closeCh closeCh

	session *session
	crypto  bool
	writeCh writeCh
	readCh  readCh
}

func newHeartbeater(s *session, w writeCh, r readCh) *heartbeater {
	h := &heartbeater{
		wg:      &sync.WaitGroup{},
		closeCh: make(chan struct{}),

		session: s,
		crypto:  true,
		writeCh: w,
		readCh:  r,
	}

	go h.heartbeat()
	go h.receive()
	return h
}

func (h *heartbeater) heartbeat() {
	h.wg.Add(1)
	defer h.wg.Done()

	for {
		select {
		case <-h.closeCh:
			return
		case writeMsg := <-h.writeCh:
			writeMsg.errCh <- h.session.write(writeMsg.msg.id, writeMsg.msg.msg, h.crypto)
		}
	}
}

func (h *heartbeater) receive() {
	h.wg.Add(1)
	defer h.wg.Done()

	for {
		select {
		case <-h.closeCh:
			return
		default:
			id, msg, err := h.session.read(h.crypto)
			h.readCh <- &readMsg{
				msg: &message{id: id, msg: msg},
				err: err,
			}
		}
	}
}

func (h *heartbeater) close() error {
	close(h.closeCh)
	h.wg.Wait()
	return h.session.Close()
}
