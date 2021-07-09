package bonus

import (
	"fmt"
	"sync"
	"time"
)

type Bonus struct {
	peerID string
	ethAdrr string

	closeCh closeCh
	wg *sync.WaitGroup

	hbWriteCh writeCh
	hbReadCh readCh
	activeBeater *heartbeater
	//ticker *time.Ticker
}

func New(peerID, ethAdrr string) (*Bonus, error) {
	wg := &sync.WaitGroup{}

	hbWriteCh := make(chan *writeMsg)
	hbReadCh := make(chan *readMsg)

	activeSession, err := newSession(newDialer())
	if err != nil {
		return nil, err
	}

	hb := newHeartbeater(activeSession, hbWriteCh, hbReadCh)

	b := &Bonus{
		peerID: peerID,
		ethAdrr: ethAdrr,

		closeCh: make(chan struct{}),
		wg: wg,

		hbWriteCh: hbWriteCh,
		hbReadCh: hbReadCh,
		activeBeater: hb,
		//ticker: time.NewTicker(time.Second * 15),
	}

	go b.serveHeartbeater()

	return b, nil
}

func (b *Bonus) serveHeartbeater() {
	b.wg.Add(1)
	defer b.wg.Done()

	var hbErr error

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- b.closeCh:
			return
		case readMsg := <- b.hbReadCh:
			if readMsg.err != nil {
			fmt.Printf("heartbeater received error: %v\n", readMsg.err)
			hbErr = readMsg.err
			break
			}
			fmt.Printf("heartbeater received CSID %v message: %v\n", readMsg.msg.id, readMsg.msg.msg)
		case t := <- ticker.C:
			fmt.Println("Current time: ", t)
			msg := &message{
				id: CSID_ID_Heartbeat,
				msg: &Heartbeat{
					Peer:b.peerID,
					EthAddr: b.ethAdrr,
				},
			}

			errCh := make(chan error)

			writeMsg := &writeMsg{
				msg: msg,
				errCh: errCh,

			}
		 	b.hbWriteCh <- writeMsg
		 	err := <- errCh
		 	if err != nil {
				fmt.Printf("heartbeater write error: %v\n", err)
				hbErr = err
				break
			}
			fmt.Printf("heartbeater write CSID %v message: %v\n", msg.id, msg.msg)
		}

		// close current active heartbeater and rebuild a new one if any error encountered
		if hbErr != nil {
			fmt.Printf("heartbeater err: %v\n", hbErr)
			fmt.Printf("Start newHearbeater.")
			b.activeBeater.close()
			for {
				activeSession, err := newSession(newDialer())
				if err != nil {
					time.Sleep(time.Second * 5)
					continue
				}

				newHB := newHeartbeater(activeSession, b.hbWriteCh, b.hbReadCh)
				b.activeBeater = newHB
				break
			}
		}
	}
}

// Close close all bonus services in top to down order.
func (b *Bonus) Close() error  {
	close(b.closeCh)
	b.wg.Wait()
	return b.activeBeater.close()
}