tests:
	bash ./test/setup.sh
	go clean -testcache
	GIT_SSH_COMMAND="ssh -i $$(pwd)/test/tmp/id_rsa" go test -v ./test
