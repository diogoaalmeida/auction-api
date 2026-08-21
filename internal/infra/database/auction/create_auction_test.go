package auction_test

import (
	"context"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGivenAuctionDuration_WhenItElapses_ThenAuctionStatusBecomesCompleted(t *testing.T) {
	ctx := context.Background()

	container, err := mongodb.Run(ctx, "mongo:6")
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Disconnect(ctx) })

	t.Setenv("AUCTION_DURATION", "1s")
	repo := auction.NewAuctionRepository(client.Database("auctions_test"))

	newAuction, err := auction_entity.CreateAuction(
		"Vintage camera",
		"Electronics",
		"A well-kept 35mm film camera",
		auction_entity.Used)
	require.Nil(t, err)

	createErr := repo.CreateAuction(ctx, newAuction)
	require.Nil(t, createErr)

	fetched, findErr := repo.FindAuctionById(ctx, newAuction.Id)
	require.Nil(t, findErr)
	assert.Equal(t, auction_entity.Active, fetched.Status)

	time.Sleep(2 * time.Second)

	closed, findErr := repo.FindAuctionById(ctx, newAuction.Id)
	require.Nil(t, findErr)
	assert.Equal(t, auction_entity.Completed, closed.Status)
}
