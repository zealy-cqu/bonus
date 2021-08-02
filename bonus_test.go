package bonus

import (
	"fmt"
	"testing"
)

func TestBonus(t *testing.T) {
	b, err := New("peerID: testID", "ethAdrr: testEthAdrr")
	if err != nil {
		t.Errorf("Failed to create new bonus. error: %v\n", err)
	}

	_ = b
	fmt.Printf("Successful to start bonus: %#+v\n", b)
	closeCh := make(chan struct{})
	<-closeCh
}
