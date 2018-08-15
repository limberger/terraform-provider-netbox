package netbox

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/digitalocean/go-netbox/netbox"
	"github.com/digitalocean/go-netbox/netbox/client"
)


c := netbox.NewNetboxWithAPIKey("your.netbox.host:8000", "your_netbox_token")