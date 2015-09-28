package cli

import (
	"flag"
	"log"

	"github.com/doubledutch/quantum/client"
)

// SoloCli defines a command to run quantum in stand alone, or solo, mode.
// Solo mode allows the user to send a request directly to a Quantum Agent
// when they know the agent address and port.
func SoloCli(args []string) {
	fs := flag.NewFlagSet("quantum solo", flag.ExitOnError)
	fs.StringVar(&port, "p", clientPort, "Agent port")
	commonFlags(fs)
	fs.Parse(args)

	lgr := lgr()
	connConfig, err := config(poolType, certFile, keyFile, caFile, lgr)
	if err != nil {
		log.Fatal(err)
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
