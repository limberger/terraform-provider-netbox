all: test provider

test:
	go test -v ./netbox/...


reconfigure: setup_util
	@echo
	@echo !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	@echo About to destroy everything in Netbox and reconfigure now!
	@echo !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	@echo
	utils/setup -destroyAll

reconfigure_testacc: reconfigure
	go clean -testcache
	TF_ACC=1 go test -v ./netbox/... -run="TestAcc"

testacc:
	@echo
	@echo !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	@echo For tests to work, your NetBox must be configured correctly.
	@echo If tests fail due to configuration errors, you may need to either
	@echo run 'make reconfigure' and rerun 'make testacc'
	@echo -OR-
	@echo run 'make reconfigure_testacc'
	@echo !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	@echo
	go clean -testcache
	TF_ACC=1 go test -v ./netbox/... -run="TestAcc"

provider:
	go build -o build/terraform-provider-netbox

# release: release_bump release_build

# release_bump:
# 	scripts/release_bump.sh

# release_build:
# 	scripts/release_build.sh

clean:
	rm -rf build/
	rm -f utils/setup

setup_util:
	cd utils; go build setup.go
