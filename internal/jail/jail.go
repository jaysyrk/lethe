package jail

import (
	"log"
	"sync"
	"time"
)

var (
	bannedIPs = make(map[string]time.Time)
	mu        sync.RWMutex
	banTime   = 24 * time.Hour
)

func BanIP(ip string) {
	mu.Lock()
	defer mu.Unlock()
	bannedIPs[ip] = time.Now().Add(banTime)
	log.Printf("IP %s has been JAILED for 24 hours", ip)
}

func IsBanned(ip string) bool {
	mu.RLock()
	defer mu.RUnlock()

	expiry, exists := bannedIPs[ip]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		return false
	}

	return true
}
