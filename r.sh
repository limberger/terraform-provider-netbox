#!/bin/bash
make build
cp pkg/linux_amd64/terraform-provider-netbox $GOPATH/bin
terraform init
terraform plan
terraform apply
