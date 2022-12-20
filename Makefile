default: help

.PHONY: help
help: # show help for each of the Makefile command
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: go-fmt 
go-fmt: # format go files with goimports
	goimports -l -w src/

.PHONY: go-tidy 
go-tidy: # run go mod tidy for each go module
	@./scripts/go-tidy.sh

.PHONY: go-test 
go-test: # run go unit tests each go module
	@./scripts/go-test.sh

.PHONY: run-docker 
run-docker: # run web services in docker
	docker compose up -d --build