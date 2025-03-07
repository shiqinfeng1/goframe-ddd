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
# gen.run: gen.wire gen.pb gen.clean 
gen.run: gen.clean gen.pb gen.ent

.PHONY: gen.pb
gen.pb: tools.verify.gf tools.verify.protoc tools.verify.protoc-gen-go tools.verify.protoc-gen-go-grpc
	@echo "===========> Generating pb files *.go from proto file through gf"
	@gf gen pb -p ${ROOT_DIR}/api/grpc -a ${ROOT_DIR}/api/grpc -c ${ROOT_DIR}/internal/server/grpc
	@gf gen ctrl -s ${ROOT_DIR}/api/http -d ${ROOT_DIR}/internal/server/http -m 

.PHONY: gen.wire
gen.wire: tools.verify.wire
	@echo "===========> Generating wire_gen.go from wire.go file through wire"
	@find internal  -mindepth 2 -maxdepth 2 | grep server | $(wireCmd)

.PHONY: gen.clean
gen.clean:
# @rm -rf ./api/client/{clientset,informers,listers}
	@find ${ROOT_DIR}/api/grpc -type f -name '*.go' -delete

.PHONY: gen.ent
gen.ent:
	@go generate ./internal/adapters/ent
