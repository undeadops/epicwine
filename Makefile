help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build: Build epicwine binary
build:
	go mod download
	go build -o epicwine cmd/epicwine/*.go

## run: Build and Run epicwine with default options
run:
	go mod download
	go run cmd/epicwine/*.go

## docker: Build a Docker Container taged as epicwine:latest
docker:
	curl -o wine.csv -s https://media.githubusercontent.com/media/m-clark/noiris/master/data-raw/wine/winemag-data-130k-v2.csv
	docker build -t epicwine:latest .

## getwine: Download wine csv file, using curl.  Creates ./wine.csv
getwine:
	curl -o wine.csv -s https://media.githubusercontent.com/media/m-clark/noiris/master/data-raw/wine/winemag-data-130k-v2.csv
