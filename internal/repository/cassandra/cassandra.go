package cassandra

import (
	"strconv"

	"github.com/gocql/gocql"
)

type CassandraDB struct {
	Session *gocql.Session
}

func NewDB(hosts []string, username, password, port, keyspace string) (*CassandraDB, error) {
	cluster := gocql.NewCluster(hosts...)
	// cluster.Keyspace = keyspace
	// cluster.Consistency = gocql.Quorum
	if port != "" {
		cluster.Port, _ = strconv.Atoi(port)
	}
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &CassandraDB{Session: session}, nil
}
