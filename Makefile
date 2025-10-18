dockerUp:
	docker compose up -d

dockerDown:
	docker compose down

fmt:
	go fmt ./...

vet:
	go vet ./...

http: fmt vet
	go run . http

migrateUp:
	go run . migrateUp

migrateDown:
	go run . migrateDown