package auth

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLoginLimiterBoundsMemoryAndRecoversAfterWindow(t *testing.T) {
	now := time.Now()
	limiter := NewLoginLimiter()
	limiter.maxEntries = 2
	limiter.now = func() time.Time { return now }
	for i := 0; i < 20; i++ {
		limiter.RecordFailure(fmt.Sprint(i), fmt.Sprint(i))
	}
	if len(limiter.ipFailures) != 2 || len(limiter.credentialFailures) != 2 {
		t.Fatal("limiter exceeded key capacity")
	}
	if allowed, retry := limiter.Check("new-ip", "new-email"); allowed || retry <= 0 {
		t.Fatal("new key allowed at capacity")
	}
	if allowed, _ := limiter.Check("0", "0"); !allowed {
		t.Fatal("known key incorrectly rejected before its failure threshold")
	}
	now = now.Add(loginFailureWindow)
	if allowed, _ := limiter.Check("new-ip", "new-email"); !allowed {
		t.Fatal("expired capacity did not recover")
	}
}

func TestLoginLimiterUsesFixedSizeNormalisedKeys(t *testing.T) {
	if len(normalizeLoginLimiterValue(strings.Repeat("x", 100000))) != 64 {
		t.Fatal("unbounded map key")
	}
	if normalizeLoginLimiterValue(" ADMIN@example.com ") != normalizeLoginLimiterValue("admin@example.com") {
		t.Fatal("normalisation changed")
	}
}
