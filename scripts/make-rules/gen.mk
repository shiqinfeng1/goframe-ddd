# ==============================================================================
# Makefile helper functions for generate necessary files
#

# generate 
ifeq ($(GOOS), darwin)
	wireCmd=xargs -I F sh -c 'cd F && echo && wire'
else
	wireCmd=xargs -i sh -c 'cd {} && echo && wire'
endif

.PHONY: gen.run
gen.run: gen.clean gen.pb gen.ctrl gen.generate gen.wire gen.mock

.PHONY: gen.pb
gen.pb: tools.verify.gf tools.verify.protoc tools.verify.protoc-gen-go tools.verify.protoc-gen-go-grpc
	@echo "===========> Generating pb files *.go from proto file through gf"
# @gf gen pb -p ${ROOT_DIR}/api/grpc -a ${ROOT_DIR}/api/grpc -c ${ROOT_DIR}/internal/server/grpc
	@for dir in $(shell find $(ROOT_DIR)/api -mindepth 1 -maxdepth 1 -type d); do \
        name=$$(basename $$dir); \
		if [ -d "$(ROOT_DIR)/api/$$name/grpc" ]; then \
			mkdir -p ${ROOT_DIR}/internal/$$name/server/grpc; \
        	gf gen pb -p $(ROOT_DIR)/api/$$name/grpc -a ${ROOT_DIR}/api/$$name/grpc -c ${ROOT_DIR}/internal/$$name/server/grpc; \
		fi \
    done 

.PHONY: gen.ctrl
gen.ctrl: tools.verify.gf
	@echo "===========> Generating ctrl files *.go from api file through gf"
# @gf gen ctrl -s ${ROOT_DIR}/api/http -d ${ROOT_DIR}/internal/server/http -m 
	@for dir in $(shell find $(ROOT_DIR)/api -mindepth 1 -maxdepth 1 -type d); do \
        name=$$(basename $$dir); \
		if [ -d "$(ROOT_DIR)/api/$$name/http" ]; then \
    		gf gen ctrl -s $(ROOT_DIR)/api/$$name/http -d $(ROOT_DIR)/internal/$$name/server/http -m; \
		fi \
    done 

.PHONY: gen.wire
gen.wire: tools.verify.wire
	@echo "===========> Generating wire_gen.go from wire.go file through wire"
	@find cmd  -mindepth 1 -maxdepth 1  | $(wireCmd)

.PHONY: gen.clean
gen.clean:
# @rm -rf ./api/client/{clientset,informers,listers}
	
	@for dir in $(shell find $(ROOT_DIR)/api -mindepth 1 -maxdepth 1 -type d); do \
        name=$$(basename $$dir); \
		if [ -d "$(ROOT_DIR)/api/$$name/grpc" ]; then \
        	find ${ROOT_DIR}/api/$$name/grpc -type f -name '*.go' -delete; \
		fi \
    done 

.PHONY: gen.generate
gen.generate:
	@echo "===========> Running go generate ./..."
	@go generate ./...

.PHONY: gen.mock
gen.mock: tools.verify.mockery
	@echo "===========> Running mockery"
	@mockery