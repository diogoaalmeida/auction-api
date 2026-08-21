# auction-api

Auction API: create auctions, place bids, find winners. Auctions close automatically once their configured duration elapses.

## Run

```bash
docker compose up --build
```

Starts the API on `:8080` and MongoDB. Env vars come from `cmd/auction/.env`.

## Configuring auction duration

| Variable | Meaning | Default |
|---|---|---|
| `AUCTION_DURATION` | How long an auction stays open before auto-closing | 5 minutes |

Go duration format, e.g. `20s`, `5m`. `AUCTION_INTERVAL` is separate: it only rejects late bids in-memory and doesn't close the auction.

## Try it

```bash
curl -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{"product_name":"Vintage camera","category":"Electronics","description":"A well-kept 35mm film camera","condition":1}'

curl http://localhost:8080/auction/<auctionId>
```

`status` starts at `0` (open) and flips to `1` (closed) on its own once `AUCTION_DURATION` passes.

## How it works

`create_auction.go` starts a goroutine after each insert that sleeps for `AUCTION_DURATION`, then updates the status to closed. It uses its own background context, since it outlives the request that created the auction.

## Run tests

```bash
go test ./...
```

Uses testcontainers-go to spin up a real MongoDB and prove the auto-close happens. Needs Docker.
