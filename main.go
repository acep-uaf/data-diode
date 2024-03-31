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

	analysis "github.com/acep-uaf/data-diode/insights"
	utility "github.com/acep-uaf/data-diode/utility"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var (
	SemVer         string
	BuildInfo      string
	ConfigSettings = "config/settings.yaml"
	InputTextFile  = "docs/example.txt"
)

type Configuration struct {
	Input struct {
		IP      string
		Port    int
		Timeout int
	}
	Output struct {
		IP   string
		Port int
		TLS  bool
	}
	Broker struct {
		Server string
		Port   int
		Topic  string
	}
}

func main() {
	data, err := os.ReadFile(ConfigSettings)

	if err != nil {
		panic(err)
	}

	var config Configuration

	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	// Configuration Settings

	diodeInputSideIP := config.Input.IP
	diodePassthroughPort := config.Input.Port
	clientLocation := fmt.Sprintf("%s:%d", diodeInputSideIP, diodePassthroughPort)

	targetServerIP := config.Output.IP
	targetServerPort := config.Output.Port
	serverLocation := fmt.Sprintf("%s:%d", targetServerIP, targetServerPort)

	mqttBrokerIP := config.Broker.Server
	mqttBrokerPort := config.Broker.Port
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
					fmt.Println(">> Client IP: ", diodeInputSideIP)
					fmt.Println(">> Client Port: ", diodePassthroughPort)
					utility.StartPlaceholderClient(diodeInputSideIP, diodePassthroughPort)
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Output side of the data diode",
				Action: func(sCtx *cli.Context) error {
					fmt.Println("----- OUTPUT -----")
					fmt.Println(">> Server IP: ", targetServerIP)
					fmt.Println(">> Server Port: ", targetServerPort)
					utility.StartPlaceholderServer(targetServerIP, targetServerPort)
					return nil
				},
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "Testing state synchronization via diode I/O",
				Action: func(tCtx *cli.Context) error {
					fmt.Println("----- TEST -----")
					analysis.Pong()
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
					return nil
				},
			},
			{
				Name:    "mqtt-subscribe",
				Aliases: []string{"ms"},
				Usage:   "Recieve payload, encapsulate payload message, & stream to diode",
				Action: func(msCtx *cli.Context) error {
					utility.InboundMessageFlow(mqttBrokerIP, mqttBrokerPort, mqttBrokerTopic, clientLocation)
					return nil
				},
			},
			{
				Name:    "mqtt-publish",
				Aliases: []string{"mp"},
				Usage:   "Recieve stream (stdin), detect complete message, decode the payload, & publish the payload",
				Action: func(mpCtx *cli.Context) error {
					utility.OutboundMessageFlow(mqttBrokerIP, mqttBrokerPort, mqttBrokerTopic, serverLocation)
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the version of the diode CLI",
				Action: func(vCtx *cli.Context) error {
					fmt.Println(">> diode version:", SemVer)
					fmt.Println(">> build information: ", BuildInfo)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal()
	}
}
