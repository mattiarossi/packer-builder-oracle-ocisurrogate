package main

import (
	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/mattiarossi/packer-builder-oracle-ocisurrogate/pkg/ocisurrogate"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(ocisurrogate.Builder))
	server.Serve()
}
