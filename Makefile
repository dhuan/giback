all: build

build:
	go build -o giback
	mkdir -p bin
	mv ./giback ./bin/.

tests:
	bash ./scripts/run_tests.sh
