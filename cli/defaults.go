package cli

import "os"

var (
	agent       string
	port        string
	requestType string
	requestData string
	logLevels   string

	serverAddr    string
	defaultOutput = os.Stdout
)

const (
	defaultHost        = "127.0.0.1"
	defaultRequestType = "basic"
	clientPort         = ":8814"
	serverPort         = ":8818"
	defaultLogLevel    = "E"
)
