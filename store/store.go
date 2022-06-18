package store

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go-inventory/objects"
)

// IStockStore is the database interface for storing Stocks
type IStockStore interface {
	Get(ctx context.Context, in *objects.GetRequest) (*objects.Stock, error)
	List(ctx context.Context, in *objects.ListRequest) ([]*objects.Stock, error)
	Create(ctx context.Context, in *objects.CreateRequest) error
	UpdateDetails(ctx context.Context, in *objects.UpdateDetailsRequest) error
}

func init() {
	rand.Seed(time.Now().UTC().Unix())
}

// GenerateUniqueID will returns a time based sortable unique id
func GenerateUniqueID() string {
	word := []byte("0987654321")
	rand.Shuffle(len(word), func(i, j int) {
		word[i], word[j] = word[j], word[i]
	})
	now := time.Now().UTC()
	return fmt.Sprintf("%010v-%010v-%s", now.Unix(), now.Nanosecond(), string(word))
}
