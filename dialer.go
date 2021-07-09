package bonus

import (
	"net"
	"sync"
	"time"
)

var testAdrr = "192.168.101.134:9527"

var goodAddrs = map[string]goodAddr{}
var badAddrs = map[string]badAddr{}

type dialer struct {
	totalCount int
	rspCount int
	wg *sync.WaitGroup
	//closeCh chan struct{}
	rspGoodCh chan goodAddr
	rspBadCh chan badAddr
	firstrspCh chan net.Conn
	failDailCh chan struct{}
}

type goodAddr struct {
	addr string
	port string
	conn net.Conn
}

type badAddr struct {
	addr string
	port string
	err error
}

func newDialer() *dialer {
	return &dialer{}
}

func (d *dialer) conn() (net.Conn, error) {
	//for addr, badAddr :=  range badAddrs {
	//	goodAddrs[addr] = goodAddr{badAddr.addr, badAddr.port, nil}
	//}
	//
	//for addr, goodAdrr := range goodAddrs {
	//	go d.dail_(addr, goodAdrr.port)
	//}
	//
	//select {
	//case conn := <- d.firstrspCh:
	//	return conn, nil
	//case <- d.failDailCh:
	//	return nil, errors.New("Failed to get connected to any server")
	//}

	return net.DialTimeout("tcp",testAdrr, time.Second*10)
}

func (d *dialer) dail_(addr, port string)  {
	d.wg.Add(1)
	defer d.wg.Done()

	conn, err := net.Dial("tcp", addr + ":" + port)
	if err != nil {
		d.rspBadCh <- badAddr{addr, port, err}
	}
	d.rspGoodCh <- goodAddr{addr, port, conn}
}

func (d *dialer)  deal() {
	d.wg.Add(1)
	defer d.wg.Done()

	var goodCount = 0
	for d.rspCount < d.totalCount {
		select {
		case goodRsp := <- d.rspGoodCh:
			if goodCount == 0 {
				d.firstrspCh <- goodRsp.conn
				goodCount++
			}

			_ = goodRsp.conn.Close()
			d.rspCount++
		case badRsp := <- d.rspBadCh:
			badAddrs[badRsp.addr] = badRsp
			delete(goodAddrs, badRsp.addr)
			d.rspCount++
		}
	}

	if goodCount == 0 {
		close(d.failDailCh)
	}
}