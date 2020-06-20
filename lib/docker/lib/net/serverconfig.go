package net



const (
	// DefaultServerPort is the default listening port for the  server
	DefaultServerPort = 8000

	// DefaultServerAddr is the default listening address for the  server
	DefaultServerAddr = "0.0.0.0"
)


type ServerConfig struct {
	// Listening port for the server
	Port int `def:"8000" opt:"p" help:"Listening port of server"`
	// Bind address for the server
	Address string `def:"0.0.0.0" help:"Listening address of server"`
	// Enables debug logging
	Debug bool `def:"false" opt:"d" help:"Enable debug level logging"`

}
