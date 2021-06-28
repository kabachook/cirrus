/*
Constant keys for viper config
*/
package cmd

import (
	"github.com/kabachook/cirrus/pkg/provider/gcp"
	"github.com/kabachook/cirrus/pkg/provider/yc"
)

const Gcp = gcp.Name
const GcpProject = Gcp + ".project"
const GcpKey = Gcp + ".key"
const GcpZones = Gcp + ".zones"

const Yc = yc.Name
const YcFolderId = Yc + ".folderId"
const YcToken = Yc + ".token"
const YcZones = Yc + ".zones"

const Server = "server"
const ServerListen = Server + ".listen"
const ServerProviders = Server + ".providers"
