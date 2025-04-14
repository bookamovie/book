# CMD #####

cmd_book := cmd/book/main.go
cmd_migrator := cmd/migrator/main.go

run_local: $(cmd_book)
	CONFIG_PATH=config/local.yaml LOG_MODE=local go run $(cmd_book)

run_dev: $(cmd_book)
	CONFIG_PATH=config/dev.yaml LOG_MODE=dev go run $(cmd_book)

run_prod: $(cmd_book)
	CONFIG_PATH=config/prod.yaml LOG_MODE=prod go run $(cmd_book)

run_custom: $(cmd_book)
	CONFIG_PATH=config/custom.yaml LOG_MODE=custom go run $(cmd_book)

migrate: $(cmd_migrator)
	MIGRATIONS=migrations/sqlite DATABASE=storage/db.sqlite go run $(cmd_migrator)