package bonus

import "testing"

func TestDialerPing(t *testing.T) {
	d := newDialer()
	d.ping()
	t.Logf("pongIPMp len: %v\n", len(d.pongIPMp))
}

func TestDoPing(t *testing.T)  {
	recv := pongIPMp{}
	gg := "www.google.com"
	bd := "www.baidu.com"
	lc := "127.0.0.1"

	err := ping_(lc, recv)
	if err != nil {
		t.Errorf("ERROR: %v\n", err)
	}

	err = ping_(bd, recv)
	if err != nil {
		t.Errorf("ERROR: %v\n", err)
	}


	err = ping_(gg, recv)
	if err != nil {
		t.Errorf("ERROR: %v\n", err)
	}
	t.Logf("\nrecv: %v\n", recv)
}