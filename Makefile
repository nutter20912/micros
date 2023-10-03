
service:=
name:=
path:=

# cannot in docker
.PHONY: migrate-diff
migrate-diff:
	atlas migrate diff --env gorm $(name).up --var service=$(service)

.PHONY: migrate-hash
migrate-hash:
	atlas migrate hash --env gorm --var service=$(service)

.PHONY: debug
debug::
ifeq ($(strip $(path)),)
	$(error 請輸入 path)
endif
debug::
	dlv debug $(path) --headless --listen=:12345 --api-version=2

.PHONY: submodule
submodule:
	git submodule update --init --recursive --remote
