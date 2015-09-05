package cli

import (
	"flag"
	"log"
	"time"

	"github.com/doubledutch/lager"
	"github.com/doubledutch/mux/gob"
	"github.com/doubledutch/quantum"
	"github.com/doubledutch/quantum/client"
)

// SoloCli defines a command to run quantum in stand alone, or solo, mode.
// Solo mode allows the user to send a request directly to a Quantum Agent
// when they know the agent address and port.
func SoloCli(args []string) {
	fs := flag.NewFlagSet("quantum solo", flag.ExitOnError)
	fs.StringVar(&agent, "agent", "", "Resolvable name or address of agent")
	fs.StringVar(&port, "p", clientPort, "Agent port")
	fs.StringVar(&requestType, "t", defaultRequestType, "Type of request")
	fs.StringVar(&requestData, "d", "{}", "Request data json")
	fs.StringVar(&logLevels, "log", defaultLogLevel, "Log levels")
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

	connConfig := &quantum.ConnConfig{
		Timeout: 100 * time.Millisecond,
		Config:  config,
	}

	log.Println("Running quantum solo")

	client := client.New(connConfig)
	conn, err := client.Dial(agent + port)
	if err != nil {
		log.Fatal(err)
	}

	request, err := NewRequest(requestType, requestData)
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}
	err = Run(conn, request, defaultOutput)
	if err != nil {
		log.Fatalf("quantum solo exited with error:\n\t%v\n", err.Error())
	} else {
		log.Println("quantum solo exited")
	}
}
