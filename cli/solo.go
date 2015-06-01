package cli

import (
	"flag"
	"log"

	"github.com/doubledutch/lager"
	"github.com/doubledutch/mux/gob"
	"github.com/doubledutch/quantum"
	"github.com/doubledutch/quantum/client"
)

// SoloCli defines a command to run quantum in stand alone, or solo, mode.
// Solo mode allows the client to send a request directly to an agent
// if you know the agent addr and port.
func SoloCli(args []string) {
	fs := flag.NewFlagSet("quantum solo", flag.ExitOnError)
	fs.StringVar(&agent, "agent", "", "name of agent")
	fs.StringVar(&port, "p", clientPort, "agent port")
	fs.StringVar(&requestType, "t", defaultRequestType, "type of request")
	fs.StringVar(&requestData, "d", "{}", "request data json")
	fs.StringVar(&logLevels, "log", defaultLogLevel, "log level")
	fs.Parse(args)

	// We need to move this
	lager := lager.NewLogLager(&lager.LogConfig{
		Levels: lager.LevelsFromString(logLevels),
		Output: defaultOutput,
	})

	config := &quantum.Config{
		Lager: lager,
		Pool:  new(gob.Pool),
	}

	log.Println("Running quantum solo")

	client := client.New(&client.Config{
		Config: config,
	})
	conn, err := client.Dial(agent + port)
	if err != nil {
		log.Fatal(err)
	}

	err = Run(conn, quantum.NewRequest(requestType, requestData), defaultOutput)
	if err != nil {
		log.Fatalf("quantum solo exited with error:\n\t%v\n", err.Error())
	} else {
		log.Println("quantum solo exited")
	}
}
