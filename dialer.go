package bonus

import (
	"errors"
	"fmt"
	"github.com/go-ping/ping"
	"math/rand"
	"net"
	"sync"
	"time"
)

var addrs = ipPortMp{
	"192.168.101.134": "9527",
	"202.182.98.210":  "1398",
	"127.0.0.1":       "44199",
}

type ipPortMp map[string]string

type pongIPMp map[int]string

type notPongIPMp map[int]string

type pongRes struct {
	ip  string
	err error
}

type dialer struct {
	pingTimeout   time.Duration
	dialTimeout   time.Duration
	retryInterval time.Duration

	targets   ipPortMp
	pingIPCh  chan string
	pongResCh chan *pongRes
	pongIPMp  pongIPMp

	pingwg *sync.WaitGroup
	pongwg *sync.WaitGroup
}

func newDialer() *dialer {
	return &dialer{
		pingTimeout:   time.Second * 10,
		dialTimeout:   time.Second * 10,
		retryInterval: time.Second * 10,

		targets:   addrs,
		pingIPCh:  make(chan string, len(addrs)),
		pongResCh: make(chan *pongRes, len(addrs)),
		pongIPMp:  make(pongIPMp, len(addrs)),

		pingwg: &sync.WaitGroup{},
		pongwg: &sync.WaitGroup{},
	}
}

// conn
func (d *dialer) conn() (net.Conn, error) {
	for {
		d.ping()

		conn, err := d.dial()
		if err != nil {
			fmt.Printf("Failed to dail %v at first round \n and will retry after %v seconds.\n", addrs, d.retryInterval)
			time.Sleep(d.retryInterval)
			continue
		}

		return conn, nil

	}

	return nil, errors.New("Failed to establish a connection")
}

// dial recursively and randomly dials from provided pongIPMp until first ip connected successfully.
func (d *dialer) dial() (net.Conn, error) {
	if len(d.pongIPMp) < 1 {
		return nil, nil
	}

	rand.Seed(86)

	n := rand.Intn(len(d.pongIPMp))
	conn, err := net.DialTimeout("tcp", d.pongIPMp[n], d.dialTimeout)
	if err != nil {
		fmt.Printf("Failed to dial %v \n", d.pongIPMp[n])
		delete(d.pongIPMp, n)
		d.dial()
	}
	fmt.Printf("connected to remote: %v\n", conn.RemoteAddr())
	return conn, nil
}

// ping pings provided ips and returns effective one,
// which should ping successfully with 10 seconds.
func (d *dialer) ping() {
	// send ip for ping through channel
	go func() {
		d.pongwg.Add(1)
		defer d.pongwg.Done()

		for ip, _ := range d.targets {
			d.pingIPCh <- ip
		}
		close(d.pingIPCh)
	}()

	// ping subroutines
	go func() {
		d.pongwg.Add(1)
		defer d.pongwg.Done()

		for ip := range d.pingIPCh {
			fmt.Printf("Createed pinger on %v\n\n", ip)
			p, err := ping.NewPinger(ip)
			if err != nil {
				fmt.Printf("Failed to create pinger on :%v\n", ip)
				return
			}

			p.Count = 1
			p.Timeout = d.pingTimeout

			go func() {
				d.pingwg.Add(1)
				defer d.pingwg.Done()

				defer p.Stop()
				err = p.Run()
				d.pongResCh <- &pongRes{ip: p.Addr(), err: err}
			}()
		}
		d.pingwg.Wait()
		//close(d.pongResCh)
	}()

	// receive subroutine
	go func() {
		d.pongwg.Add(1)
		defer d.pongwg.Done()

		var pongN = -1
		var pongSucN = -1
		for pongRes := range d.pongResCh {
			pongN++
			if pongRes.err == nil {
				pongSucN++
				d.pongIPMp[pongN] = pongRes.ip
			}

			if pongN == len(d.targets) {
				close(d.pongResCh)
				return
			}
		}
	}()

	d.pongwg.Wait()
}

func pings(ips ipPortMp) (pongIPMp, notPongIPMp) {
	return nil, nil
}

// todo: use concurrent map
func ping_(ip string, recv pongIPMp) error {
	p, err := ping.NewPinger(ip)
	if err != nil {
		fmt.Printf("Failed to create pinger for :%v\n", ip)
		return err
	}

	p.Timeout = time.Second * 10

	p.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		p.Stop()
	}

	p.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

		if stats.PacketsRecv >= 1 {
			recv[len(recv)+1] = stats.Addr
		}
	}

	return p.Run()
}
