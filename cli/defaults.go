package cli

import (
	"errors"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/doubledutch/lager"
	"github.com/doubledutch/mux"
	"github.com/doubledutch/mux/gob"
	"github.com/doubledutch/quantum"
	"github.com/doubledutch/quantum/client"
)

var (
	server      string
	agent       string
	port        string
	requestType string
	requestData string
	poolType    string
	logLevels   string
	certFile    string
	keyFile     string
	caFile      string
	tlsVerify   bool

	defaultOutput = os.Stdout
)

const (
	defaultHost        = "127.0.0.1"
	defaultRequestType = "basic"
	clientPort         = ":8814"
	serverPort         = ":8818"
	defaultLogLevel    = "E"
	defaultPool        = "gob"
)

// commonFlags configurse FlagSet to read common flags
func commonFlags(fs *flag.FlagSet) {
	fs.StringVar(&agent, "agent", defaultHost, "Agent name to resolve")
	fs.StringVar(&requestType, "t", defaultRequestType, "Type of request")
	fs.StringVar(&requestData, "d", "{}", "Request data json")
	fs.StringVar(&logLevels, "log", defaultLogLevel, "Log levels")
	fs.StringVar(&poolType, "pool", defaultPool, "mux pool type")
	fs.StringVar(&certFile, "tlscert", "", "certificate file path")
	fs.StringVar(&keyFile, "tlskey", "", "certificate key path")
	fs.StringVar(&caFile, "tlsca", "", "certificate CA path")
	fs.BoolVar(&tlsVerify, "tlsverify", false, "verify the server's certificate chain and host name")
}

func lgr() lager.Lager {
	return lager.NewLogLager(&lager.LogConfig{
		Levels: lager.LevelsFromString(logLevels),
		Output: defaultOutput,
	})
}

func config(poolType, certFile, keyFile, caFile string, lgr lager.Lager) (*quantum.ConnConfig, error) {
	pool := muxPool(poolType)
	if pool == nil {
		return nil, errors.New("Invalid pool provided")
	}

	cc := &quantum.ConnConfig{
		Timeout: 100 * time.Millisecond,
		Config: &quantum.Config{
			Lager: lgr,
			Pool:  pool,
		},
	}

	if certFile != "" && keyFile != "" && caFile != "" {
		var err error
		cc.TLSConfig, err = client.NewTLSConfig(certFile, keyFile, caFile)
		if err != nil {
			return nil, err
		}
		cc.TLSConfig.InsecureSkipVerify = !tlsVerify
		cc.TLSConfig.ServerName = agent

		// Increase timeout for TLS handshake
		cc.Timeout = 10 * time.Second
	}

	return cc, nil
}

func muxPool(str string) mux.Pool {
	switch strings.ToLower(str) {
	case "gob":
		return new(gob.Pool)
	}

	return nil
}
