gen_auth_db_mocks:
	mockgen -source=services/auth/dataservice/dataservice.go \
	-destination=services/auth/dataservice/mocks/mock_dataservice.go

gen_user_db_mocks:
	mockgen -source=services/user/dataservice/dataservice.go \
	-destination=services/user/dataservice/mocks/mock_dataservice.go

gen_ai_db_mocks:
	mockgen -source=services/ai/dataservice/dataservice.go \
	-destination=services/ai/dataservice/mocks/mock_dataservice.go

gen_auth_adapter_mocks:
	mockgen -source=services/auth/adapter/adapter.go \
	-destination=services/auth/adapter/mocks/mock_adapter.go

coverage:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out	