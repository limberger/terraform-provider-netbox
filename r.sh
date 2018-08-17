#!/bin/bash
# Run a sample
make build
cp pkg/linux_amd64/terraform-provider-netbox $GOPATH/bin
mkdir tmp
export TF_TF_LOG_PATH=tmp/log
export TF_LOG=DEBUG
terraform init
terraform plan
terraform apply
