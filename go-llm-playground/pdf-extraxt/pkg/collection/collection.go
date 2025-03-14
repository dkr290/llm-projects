package collection

import (
	"context"

	"github.com/qdrant/go-client/qdrant"
)

func uint32Ptr(i uint32) *uint32 { return &i }

func CreateCollection(coll string, host string) error {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		Port:   6334,
		UseTLS: false,
	})
	if err != nil {
		return err
	}
	err = client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName:    coll,
		ShardNumber:       uint32Ptr(6),
		ReplicationFactor: uint32Ptr(3),
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     768,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		return err
	}
	return nil
}
