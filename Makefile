all: build

build:
	go build -o giback
	mkdir -p bin
	mv ./giback ./bin/.

tests:
	bash ./scripts/run_tests.sh

docs_build: 
	bash ./scripts/docs_build.sh
