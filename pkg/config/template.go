package config

// TOML service template

// Config represents the global micro-service configuration.
type Config struct {
	ConfigFile string `toml:"-"`
	// List of services
	Alfa  Service `toml:"alfa"`
	Beta  Service `toml:"beta"`
	Omega Service `toml:"omega"`
}

// Service represents a micro-service struct.
type Service struct {
	Name      string
	ShortName string   `toml:"short_name" json:",omitempty"`
	Database  Database `toml:"database"`
	Server    Server
}

// Server represents a server configuration struct.
type Server struct {
	Host          string
	Port          int
	Address       string
	Clients       []string `json:",omitempty"`
	MicroServices []string `toml:"micro_services" json:",omitempty"`
	IsGateway     bool     `toml:"is_gateway" json:",omitempty"`
	GatewayPort   int      `toml:"gateway_port" json:",omitempty"`
}

// Database represents a database struct.
type Database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
}
