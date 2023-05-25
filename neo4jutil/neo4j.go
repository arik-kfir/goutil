package neo4jutil

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
)

type Neo4jZerologBoltLogger struct {
	Logger *zerolog.Logger
}

func (nl *Neo4jZerologBoltLogger) LogClientMessage(id, msg string, args ...interface{}) {
	nl.Logger.Trace().
		Str("type", "client").
		Str("id", id).
		Msgf(msg, args...)
}

func (nl *Neo4jZerologBoltLogger) LogServerMessage(id, msg string, args ...interface{}) {
	nl.Logger.Trace().
		Str("type", "server").
		Str("id", id).
		Msgf(msg, args...)
}

func (c *Neo4jConfig) Connect() (neo4j.DriverWithContext, error) {
	var url string
	if c.TLS {
		url = fmt.Sprintf("bolt://%s:%d", c.Host, c.Port)
	} else {
		url = fmt.Sprintf("bolt+s://%s:%d", c.Host, c.Port)
	}
	neo4jDriver, err := neo4j.NewDriverWithContext(url, neo4j.NoAuth())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Neo4j: %w", err)
	}
	return neo4jDriver, nil
}
