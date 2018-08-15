package main

import (
	"github.com/hashicorp/terraform/plugin"
	"gothub.com/limberger"
)

func main() {
	plugin.Serve(new(MyPlugin))
}
