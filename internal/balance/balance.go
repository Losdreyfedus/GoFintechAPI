package balance

import (
	"encoding/json"
	"sync"
	"time"
)

type Balance struct {
	UserID        int       `json:"user_id"`
	Amount        float64   `json:"amount"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	mu            sync.RWMutex
}

func (b *Balance) Update(amount float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Amount = amount
	b.LastUpdatedAt = time.Now()
}

func (b *Balance) GetAmount() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Amount
}

func (b *Balance) MarshalJSON() ([]byte, error) {
	type Alias Balance
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	})
}
