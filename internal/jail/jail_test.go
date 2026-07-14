package jail

import (
	"testing"
	"time"
)

func TestBanIP(t *testing.T) {
	ip := "192.168.1.100"


	if IsBanned(ip) {
		t.Errorf("Expected IP %s to not be banned initially", ip)
	}


	BanIP(ip)


	if !IsBanned(ip) {
		t.Errorf("Expected IP %s to be banned", ip)
	}
}

func TestIsBanned_Expired(t *testing.T) {
	ip := "10.0.0.5"


	mu.Lock()
	bannedIPs[ip] = time.Now().Add(-1 * time.Hour)
	mu.Unlock()


	if IsBanned(ip) {
		t.Errorf("Expected IP %s to not be banned (expired)", ip)
	}
}
