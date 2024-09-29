package config

import "os"

type CassandraConfig struct {
	CassandraUsername string
	CassandraPassword string
	CassandraHost     string
	CassandraPort     string
}

func CassandraLoad(env string) (*CassandraConfig, error) {
	config := &CassandraConfig{
		CassandraUsername: os.Getenv("CASSANDRA_USERNAME"),
		CassandraPassword: os.Getenv("CASSANDRA_PASSWORD"),
		CassandraHost:     os.Getenv("CASSANDRA_HOST"),
		CassandraPort:     os.Getenv("CASSANDRA_PORT"),
	}

	return config, nil
}
