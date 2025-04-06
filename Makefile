.PHONY: run
run:
	@go run main.go setup.go

.PHONY: generate-key
generate-key: 
	@openssl ecparam -name prime256v1 -genkey -noout -out ecdsa_private.pem
	@openssl ec -in ecdsa_private.pem -pubout -out ecdsa_public.pem

.PHONY: docker-up
docker-up: 
	@docker-compose up -d

.PHONY: docker-down
docker-down: 
	@docker-compose down

.PHONY: mock
mock: 
	@mockery