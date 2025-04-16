VERSION := 1.2.0

CONFIG_PATH ?= config/local.yaml     # ALT. CONFIGS LOCATED IN 'config' FOLDER
LOG_MODE ?= local					 # ALT. LOG MODES LOCATED IN 'intrenal/lib/logger/logger.go' FILE
MIGRATIONS ?= migrations/sqlite      # TO MIGRATE DB FOR TESTS USE 'tests/migrations/sqlite' INSTEAD
STORAGE ?= storage/db.sqlite         # FOR TESTS USE 'tests/storage/db.sqlite' INSTEAD

# CMD #####

BOOK_MAIN ?= cmd/book/main.go
MIGRATOR_MAIN ?= cmd/migrator/main.go

run: $(BOOK_MAIN)
	CONFIG_PATH=$(CONFIG_PATH) LOG_MODE=$(LOG_MODE) go run $(BOOK_MAIN)

migrate: $(MIGRATOR_MAIN)
	MIGRATIONS=$(MIGRATIONS) STORAGE=$(STORAGE) go run $(MIGRATOR_MAIN)

# TESTS ###

TYPE ?= all

test:
	@case $(TYPE) in \
		all) CONFIG_PATH=config/test.yaml LOG_MODE=$(LOG_MODE) go test ./tests -v ;; \
		functional) CONFIG_PATH=config/test.yaml LOG_MODE=$(LOG_MODE) go test ./tests -v -run _Functional ;; \
		unit) CONFIG_PATH=config/test.yaml LOG_MODE=$(LOG_MODE) go test ./tests -v -run _Unit ;; \
		integration) CONFIG_PATH=config/test.yaml LOG_MODE=$(LOG_MODE) go test ./tests -v -run _Integration ;; \
	esac

# DOCKER ##

IMAGE_NAME := book:$(VERSION)
CONTAINER_NAME := book-$(VERSION)

ACTION ?= *
EXEC ?= *

docker:
	@case $(ACTION) in \
		*) echo "Missing 'ACTION' value. specify it with 'ACTION=...'. If you trying to 'ACTION=exec', please specify the 'EXEC=...'";; \
		build) docker build -f deployments/docker/Dockerfile -t $(IMAGE_NAME) . ;; \
		run) docker run --name $(CONTAINER_NAME) -p 5092:5092 -d -e CONFIG_PATH=$(CONFIG_PATH) -e LOG_MODE=$(LOG_MODE) $(IMAGE_NAME);; \
		exec) \
			case $(EXEC) in \
				*) echo "missing 'EXEC' value. specify it with 'EXEC=...'";; \
				migrate) docker exec -it $(CONTAINER_NAME) bash -c "MIGRATIONS=$(MIGRATIONS) STORAGE=$(STORAGE) tools/migrator" ;; \
			esac ;; \
		remove) docker rm -f -v $(CONTAINER_NAME) || true; \
				docker rmi -f $(IMAGE_NAME) ;; \
		start) docker start $(CONTAINER_NAME) ;; \
		stop) docker stop $(CONTAINER_NAME) ;; \
	esac