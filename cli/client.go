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
// quantum Server. By specifing the type of the request, and optionally the
// host, quantum Server will return JobRecord(s) for available agents
// to communicate with to run the job.
func ClientCli(args []string) {
	fs := flag.NewFlagSet("quantum client", flag.ExitOnError)
	fs.StringVar(&serverAddr, "server", "", "quantum server to register with")
	fs.StringVar(&agent, "agent", "", "name of agent")
	fs.StringVar(&requestType, "t", "", "type of request")
	fs.StringVar(&requestData, "d", "{}", "request data json")
	fs.StringVar(&logLevels, "log", defaultLogLevel, "log levels")
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

	// We aren't capturing the output nor forwarding signals here.
	// cr.Resolve used to return a client, not it returns a connection.. hm
	if err = Run(conn, quantum.NewRequest(requestType, requestData), defaultOutput); err != nil {
		log.Fatalf("quantum client exited with error:\n\t%v\n", err.Error())
	} else {
		log.Println("quantum client exited")
	}
}
