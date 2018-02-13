.PHONY: help
help:
	@echo Usage: make [command]
	@echo Commands:
	@echo "	all 	- install, unpatch"
	@echo "	install 	- patch -> go install"
	@echo "	patch 	- apply all patches"
	@echo "	unpatch - undo/unapply all patches"
	@echo "	refresh - update the quilt patchset"
	@echo "	update 	- pull in new go source from GOROOT"
	@echo "	clean 	- unpatch & clean up git"
	@echo ""
	@echo "To update binary run... (* requires go workspace *)"
	@echo "	make all->clean"
	@echo ""
	@echo "To update to new source, make..."
	@echo "	make update->commit->patch->[fix]->refresh->unpatch->commit->clean"
	@echo ""
	@echo "To work on patch..."
	@echo "	make patch->[work]->refresh->unpatch->commit->clean"

all: install unpatch

install: patch
	go install

# to apply patch to work on it or build it
# noop if patch is already applied (-N -r)
patch:
	patch -p1 -N -r /dev/null < ./printer/compact.patch || true

# to remove patches for commiting them or updating upstream code
unpatch:
	git co .

# after working on a patch, to update it.
# use 'quilt add [filename]' for new files (before editing)
refresh-patch:
	git diff -u > ./printer/compact.patch

# update fmt code from GO sources
update: check-goroot
	find ${GOROOT}/src/go/printer -maxdepth 1 -type f \
		! -name '*_test.go' -exec cp {} ./printer \;

check-goroot:
ifndef GOROOT
	$(error GOROOT is undefined)
endif

clean: unpatch
	rm -f gocompact
