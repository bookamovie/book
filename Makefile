# CMD #####

cmd_book := cmd/book/main.go

run_local: $(cmd_book)
	CONFIG_PATH=config/book_local.yaml LOG_MODE=local go run $(cmd_book)

run_dev: $(cmd_book)
	CONFIG_PATH=config/book_dev.yaml LOG_MODE=dev go run $(cmd_book)

run_prod: $(cmd_book)
	CONFIG_PATH=config/book_prod.yaml LOG_MODE=prod go run $(cmd_book)

run_custom: $(cmd_book)
	CONFIG_PATH=config/book_custom.yaml LOG_MODE=custom go run $(cmd_book)

# PROTO ###

proto: 
 protoc --proto_path=protos/book --go_out=gen/go --go-grpc_out=gen/go book.proto