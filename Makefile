migrate_up:
	go run . migrate

migrate_down:
	go run . migrate -v 0

migrate_to:
	go run . migrate -v "$(v)"

test_data:
	go run . test-data -n "$(n)"

core_mocks:
	mockery --dir internal/core --all --output internal/core/mocks --note 'Regenerate with `make core_mocks`'

store_mocks:
	mockery --dir internal/storage --all --output internal/storage/mockstore --note 'Regenerate using `make store_mocks`'

.PHONY: migrate_up migrate_down migrate_to test_data core_mocks store_mocks
