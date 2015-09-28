package cli

import (
	"flag"
	"log"
	"strings"

	"github.com/doubledutch/quantum"
	"github.com/doubledutch/quantum/consul"
)

// ClientCli defines a command to run quantum in client mode, making use of
// Quantum Resolver. By specifing the type of the request, and optionally the
// hostname, Quantum Resolver will return the addresses and ports of available agents
// to communicate with to run jobs.
func ClientCli(args []string) {
	fs := flag.NewFlagSet("quantum client", flag.ExitOnError)
	fs.StringVar(&server, "server", "", "type of the server")
	commonFlags(fs)
	fs.Parse(args)

	log.Printf("Running quantum client, resolving with %s\n", server)

	lgr := lgr()
	connConfig, err := config(poolType, certFile, keyFile, caFile, lgr)
	if err != nil {
		log.Fatal(err)
	}

	var cr quantum.ClientResolver
	switch strings.ToLower(server) {
	case "consul":
		cr = consul.NewClientResolverFromEnv(connConfig)
	default:
	}

	if cr == nil {
		log.Fatalf("Invalid Client Resolver: %s", server)
	}

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
