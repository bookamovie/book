# CMD #####

cmd_book := cmd/book/main.go
cmd_migrator := cmd/migrator/main.go

run: $(cmd_book)
	go run $(cmd_book)

migrate: $(cmd_migrator)
	go run $(cmd_migrator)

# TESTS ###

TYPE ?= all

test:
	@case $(TYPE) in \
		all) go test ./tests -v ;; \
		functional) go test ./tests -v -run _Functional ;; \
		unit) go test ./tests -v -run _Unit ;; \
		integration) go test ./tests -v -run _Integration ;; \
	esac
