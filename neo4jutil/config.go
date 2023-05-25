package neo4jutil

type Neo4jConfig struct {
	Host string `env:"HOST" value-name:"HOST" long:"host" description:"Neo4j host name" required:"yes"`
	Port int    `env:"PORT" value-name:"PORT" long:"port" description:"Neo4j port" required:"yes"`
	TLS  bool   `env:"TLS" value-name:"TLS" long:"tls" description:"Whether to use TLS to connect to Neo4j"`
}
