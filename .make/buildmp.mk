GOBUILDOUTMP=-o bin/${@:build.go/%=%}$(BINARY_EXT)

TARGETS:= $(foreach P,$(BINARIES:./cmd/%=%),$(addprefix $P-,$(BUILDS)))

.PHONY=buildmp build.go/%
buildmp: tools.msign _bindir $(TARGETS:%=build.go/%) ## do project build
build.go/%: BIN = $(word 1,$(subst -, ,$*))
build.go/%: export GOOS = $(word 2,$(subst -, ,$*))
build.go/%: export GOARCH = $(word 3,$(subst -, ,$(subst ., ,$*)))
build.go/%: BINARY_EXT = $(if $(filter windows, $(word 2,$(subst -, ,$*))),.exe,$(EMPTY))
build.go/%: prebuild
ifeq ($(FEATURE_SELF_UPDATE),yes)
ifeq ($(MSIGN_PUBLIC),)
	$(error MSIGN_PUBLIC is undefined)
endif
endif
	$(GOBUILD) $(GOBUILDOUTMP) ./cmd/$(BIN)
ifeq ($(MSIGN_SIGNATURE),yes)
ifeq ($(MSIGN_PRIVATE),)
	$(error MSIGN_PRIVATE is undefined)
endif
	$(MSIGN) sign --to-file bin/${@:build.go/%=%}$(BINARY_EXT)
endif

.PHONY=cleanmp cleanmp/%
cleanmp: $(TARGETS:%=cleanmp/%) ## clean up files
	$(GOCLEAN)
	rm -f $(GOOUTDIR)/cover.out
	rm -f $(GOOUTDIR)/cover.html

cleanmp/%: BINARY_EXT = $(if $(filter windows, $(word 2,$(subst -, ,$*))),.exe,$(EMPTY))
cleanmp/%:
	rm -f ${@:cleanmp/%=$(GOOUTDIR)/%}$(BINARY_EXT)
	rm -f ${@:cleanmp/%=$(GOOUTDIR)/%}$(BINARY_EXT).msign
