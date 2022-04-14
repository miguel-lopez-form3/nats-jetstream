.PHONY: syncpub
syncpub:
	go run ./cmd/syncpub/syncpub.go

.PHONY: asyncpub
asyncpub:
	go run ./cmd/asyncpub/asyncpub.go

.PHONY: pullsub
pullsub:
	go run ./cmd/pullsub/pullsub.go

.PHONY: pushsub
pushsub:
	go run ./cmd/pushsub/pushsub.go