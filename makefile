.DEFAULT_GOAL := everything

dependencies:
	@echo Downloading Dependencies
	@go get ./...

build: dependencies
	@echo Compiling Apps
	@echo   --- keepsake 
	@go build -ldflags="-s -w" github.com/riomhaire/keepsake
	@upx  keepsake
	@cp keepsake ${GOPATH}/bin
	@echo Done Compiling Apps

buildarm: dependencies
	@echo Compiling Apps
	@echo   --- keepsake 
	@GOOS=linux GOARCH=arm GOARM=5  go build -o keepsake-arm --ldflags="-s -w" github.com/riomhaire/keepsake
	@upx  keepsake-arm
	@echo Done Compiling Apps

test:
	@echo Running Unit Tests
	@go test ./...

profile:
	@echo Profiling Code
	@go get -u github.com/haya14busa/goverage 
	@goverage -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out

clean:
	@echo Cleaning
	@go clean
	@rm -f keepsake
	@rm -f coverage*.html
	@find . -name "debug.test" -exec rm -f {} \;

everything: clean build test profile  
	@echo Done
