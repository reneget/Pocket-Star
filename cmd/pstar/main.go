package main

import (
	_ "github.com/anomalyco/my-pretty-star/modules/core"
	_ "github.com/anomalyco/my-pretty-star/modules/monitor"
	_ "github.com/anomalyco/my-pretty-star/modules/tui"
	"github.com/anomalyco/my-pretty-star/pkg/cli"
)

func main() {
	cli.Execute()
}
