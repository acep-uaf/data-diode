// ______      _         ______ _           _
// |  _  \    | |        |  _  (_)         | |
// | | | |__ _| |_ __ _  | | | |_  ___   __| | ___
// | | | / _` | __/ _` | | | | | |/ _ \ / _` |/ _ \
// | |/ / (_| | || (_| | | |/ /| | (_) | (_| |  __/
// |___/ \__,_|\__\__,_| |___/ |_|\___/ \__,_|\___|

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"rsc.io/quote"
)

func main () {
	app := &cli.App{
		Name:  "diode",
		Usage: "A command line tool for interacting with data diodes.",
		Action: func(cCtx *cli.Context) error {
			fmt.Println(quote.Go())
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Input side of the data diode",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(">> INPUT")
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Output side of the data diode",
				Action: func(sCtx *cli.Context) error {
					fmt.Println(">> OUTPUT")
					return nil
				},
			},
			{
				Name:    "diagnostics",
				Aliases: []string{"d"},
				Usage:   "Debug diagnostics via configuration settings",
				Action: func(dCtx *cli.Context) error {
					fmt.Println(">> DIAGNOSTICS")
					return nil
				},
			},
			{
				Name:    "benchmark",
				Aliases: []string{"b"},
				Usage:   "System benchmark analysis + report performance metrics",
				Action: func(bCtx *cli.Context) error {
					fmt.Println(">> BENCHMARKS")
					return nil
				},
			},
			{
				Name: "version",
				Aliases: []string{"v"},
				Usage: "Print the version of the diode CLI",
				Action: func(vCtx *cli.Context) error {
					fmt.Println(">> diode version 0.0.1")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal()
	}
}
