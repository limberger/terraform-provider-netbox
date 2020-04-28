all: test provider

build/:
	mkdir -p $@

test:
	go test -v ./netbox/...

testacc:
	@echo
	@echo ============================================================
	@echo For tests to work, your NetBox must be configured correctly.
	@echo
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
