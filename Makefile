.PHONY: help
help:
	@echo Usage: make [command]
	@echo Commands:
	@echo "	toggle-patch   - remove/apply patch"
	@echo "	refresh-patch  - update the patch"
	@echo "	update-printer - pull in new go source from GOROOT"
	@echo "	clean          - cleans up any junk (eg. .orig files)"
	@echo ""
	@echo "To update to new printer source..."
	@echo "	update-printer-> toggle-patch-> [fix->] refresh-patch-> commit"
	@echo ""
	@echo "After updating code..."
	@echo "	refresh-patch-> commit"

# reverse installed patch (or re-apply it)
toggle-patch:
	@cd printer && patch -p2 -t --posix < ./compact.patch

# update patch
refresh-patch:
	diff -u printer/orig/ printer/ \
		| grep -v '^Only in' > printer/compact.patch || true

# update fmt code from GO sources
update-printer: check-goroot
	@echo Updating from ${GOROOT}/src/go/printer/\*
	@mkdir -p printer/orig
	@find ${GOROOT}/src/go/printer -maxdepth 1 -type f \
		! -name '*_test.go' -exec cp {} ./printer/orig \;
	@cp ${GOROOT}/VERSION ./printer/
	@cp printer/orig/* printer/
	@echo "Remember to re-apply the patch (make toggle-patch)."

check-goroot:
ifndef GOROOT
	$(error GOROOT is undefined)
endif

clean:
	find . -name '*.orig' -delete
	rm -f gocompact
