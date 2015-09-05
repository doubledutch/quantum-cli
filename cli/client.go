package cli

import (
	"flag"
	"log"

	"github.com/doubledutch/lager"
	"github.com/doubledutch/mux/gob"
	"github.com/doubledutch/quantum"
	"github.com/doubledutch/quantum/consul"
)

// ClientCli defines a command to run quantum in client mode, making use of
// Quantum Resolver. By specifing the type of the request, and optionally the
// hostname, Quantum Resolver will return the addresses and ports of available agents
// to communicate with to run jobs.
func ClientCli(args []string) {
	fs := flag.NewFlagSet("quantum client", flag.ExitOnError)
	fs.StringVar(&serverAddr, "server", "", "Address of resolver")
	fs.StringVar(&agent, "agent", "", "Hostname of agent to resolve")
	fs.StringVar(&requestType, "t", "", "Type of request")
	fs.StringVar(&requestData, "d", "{}", "Request data json")
	fs.StringVar(&logLevels, "log", defaultLogLevel, "Log levels")
	fs.Parse(args)

	log.Printf("Running quantum client, resolving with %s\n", serverAddr)

	lager := lager.NewLogLager(&lager.LogConfig{
		Levels: lager.LevelsFromString(logLevels),
		Output: defaultOutput,
	})

	config := &quantum.Config{
		Lager: lager,
		Pool:  new(gob.Pool),
	}

	cr := consul.NewClientResolver(quantum.ClientResolverConfig{
		Server: serverAddr,
		Config: config,
	})
	conn, err := cr.Resolve(quantum.ResolveRequest{
		Agent: agent,
		Type:  requestType,
	})
	if err != nil {
		log.Fatal(err)
	}

	request, err := NewRequest(requestType, requestData)
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}

	if err = Run(conn, request, defaultOutput); err != nil {
		log.Fatalf("quantum client exited with error:\n\t%v\n", err.Error())
	} else {
		log.Println("quantum client exited")
	}
}
