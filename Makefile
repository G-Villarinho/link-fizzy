.PHONY: run
run:
	go run main.go setup.go

.PHONY: generate-key
generate-key: ## Generate private and public key
	@openssl ecparam -name prime256v1 -genkey -noout -out ecdsa_private.pem
	@openssl ec -in ecdsa_private.pem -pubout -out ecdsa_public.pem
