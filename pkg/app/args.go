package app

import "flag"

const (
	// defaultServerPort is the default port to listen on
	defaultServerPort = 3344
	// defaultServerAddress is the default address to listen on
	defaultServerAddress = "localhost"
)

var (
	// portFlag is the port to listen on
	portFlag = flag.Int("port", defaultServerPort, "Port to listen on")
	// serverAddressFlag is the address to listen on
	serverAddressFlag = flag.String("address", defaultServerAddress, "Address to listen on")
)

// App is the struct that holds the parsed flags
type App struct {
	Port          int
	ServerAddress string
}

// ParseFlags parses the flags and returns an App struct
func ParseFlags() *App {
	flag.Parse()

	return &App{
		Port:          *portFlag,
		ServerAddress: *serverAddressFlag,
	}
}
