VERSION=$$(git describe --tags $$(git rev-list --tags --max-count=1))

build: emr-cli
	docker build . -t bender

run:
	# eval $$(DEFAULT_get-ssm-params -path /ops/bender -output shell)
	docker run -it -e SLACK_TOKEN bender

dev:
	which justrun || go get github.com/jmhodges/justrun
	justrun -c 'go run main.go' -delay 10000ms main.go lib/**/*
