dev:
	docker-compose up -d --build
down:
	docker-compose down
protos:
	protoc -I node/proto node/proto/blockchain.proto --go_out=node --go-grpc_out=node