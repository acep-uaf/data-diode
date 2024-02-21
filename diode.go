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
	"time"

	"github.com/acep-uaf/data-diode/utility"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var (
	SemVer string
)

type Configuration struct {
	Input struct {
		IP   string
		Port int
	}
	Output struct {
		IP   string
		Port int
	}
	Broker struct {
		Server  string
		Port    int
		Topic   string
		Message string
	}
}

func sampleMetrics(server string, port int) {
	fmt.Println(">> Local time: ", time.Now())
	fmt.Println(">> UTC time: ", time.Now().UTC())
	fmt.Println(">> Value: ", utility.Value())
	// utility.Client(server, port)
}

func main() {
	data, err := os.ReadFile("config.yaml")

	if err != nil {
		panic(err)
	}

	var config Configuration

	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	// Configuration Settings

	diodeInputSideIP := config.Input.IP
	diodeTCPPassthroughPort := config.Input.Port
	targetTCPServerIP := config.Output.IP
	targetTCPServerPort := config.Output.Port

	mqttBrokerIP := config.Broker.Server
	mqttBrokerPort := config.Broker.Port
	mqttBrokerMessage := config.Broker.Message
	mqttBrokerTopic := config.Broker.Topic

	app := &cli.App{
		Name:  "diode",
		Usage: "Tool for interacting with data diode(s) via command-line interface (CLI).",
		Action: func(cCtx *cli.Context) error {
			fmt.Println("diode: try 'diode --help' for more information")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Input side of the data diode",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("----- INPUT -----")
					utility.Client(diodeInputSideIP, diodeTCPPassthroughPort)
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Output side of the data diode",
				Action: func(sCtx *cli.Context) error {
					fmt.Println("----- OUTPUT -----")
					utility.Server(targetTCPServerIP, targetTCPServerPort)
					return nil
				},
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "Testing state synchronization via diode I/O",
				Action: func(tCtx *cli.Context) error {
					fmt.Println("----- TEST -----")
					utility.Validation()
					return nil
				},
			},
			{
				Name:    "diagnostics",
				Aliases: []string{"d"},
				Usage:   "Debug diagnostics via configuration settings",
				Action: func(dCtx *cli.Context) error {
					fmt.Println("----- DIAGNOSTICS -----")
					fmt.Printf("%+v\n", config)
					return nil
				},
			},
			{
				Name:    "benchmark",
				Aliases: []string{"b"},
				Usage:   "System benchmark analysis + report performance metrics",
				Action: func(bCtx *cli.Context) error {
					fmt.Println("----- BENCHMARKS -----")
					sampleMetrics(utility.CONN_HOST, 3333)
					return nil
				},
			},
			{
				Name:    "mqtt",
				Aliases: []string{"m"},
				Usage:   "MQTT (TCP stream) demo",
				Action: func(mCtx *cli.Context) error {
					fmt.Println("----- MQTT -----")
					utility.Republisher(mqttBrokerIP, mqttBrokerPort, mqttBrokerTopic, mqttBrokerMessage)
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the version of the diode CLI",
				Action: func(vCtx *cli.Context) error {
					fmt.Println(">> diode version " + SemVer)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal()
	}
}
