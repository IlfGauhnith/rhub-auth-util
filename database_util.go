package util

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

/*
Each Function Invocation Is Independent:
Each time a Google Cloud Function is invoked, a new execution context is created.
However, this doesn't always mean a new virtual machine (VM) is created.
If there is existing capacity, an already running instance (that has been "warm") might
handle the new request, and hence, resources (such as a dbPool) can be reused across
multiple invocations within the same VM.

Instance Lifecycle:
When a Cloud Function is called for the first time, Google Cloud creates a new instance
(which could be a VM or a container). If the function is called again while this instance
is still "warm," Google will reuse that instance. This means that global variables like
dbPool will persist between invocations on the same instance.
Cold start: If no instance is available (or a function hasn't been invoked recently),
a new instance is created, and this is when your dbPool is initialized again.

Global Variables and Instance Reuse:
By declaring dbPool as a global variable, it is initialized only once per instance when the
instance starts. On subsequent invocations on the same instance, Cloud Functions won't
reinitialize it.
While each invocation is stateless, instance reuse allows the same instance
(and its connection pool) to be used across multiple invocations for efficiency.
*/

var DBPool *sql.DB
var dbErr error

func init() {
	DBPool, dbErr = ConnectToDatabase()

	if dbErr != nil {
		log.Printf("Error connecting the database: %v", dbErr)
	}
}

func ConnectToDatabase() (*sql.DB, error) {
	LogFunctionExecutionStart(ConnectToDatabase)

	// Load environment variables
	host := os.Getenv("DB_HOST_DEV")
	port := os.Getenv("DB_PORT_DEV")
	dbname := os.Getenv("DB_NAME_DEV")
	user := os.Getenv("DB_USER_DEV")
	password := os.Getenv("DB_PASSWORD_DEV")

	if host == "" || port == "" || dbname == "" || user == "" || password == "" {
		return nil, fmt.Errorf("missing environment variable for database connection")
	}

	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s",
		host, user, password, port, dbname)

	log.Printf("Connecting to database with URI: %s", dbURI)

	DBPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		log.Printf("sql.Open: %v", err)
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Set max open connections and idle connections
	DBPool.SetMaxOpenConns(5)
	DBPool.SetMaxIdleConns(3)
	DBPool.SetConnMaxIdleTime(5) // Idle time in minutes

	// Ping the database to check the connection
	err = DBPool.Ping()
	if err != nil {
		log.Printf("Error pinging the database: %v", err)
		return nil, fmt.Errorf("DB ping error: %w", err)
	}
	log.Print("Database successfully pinged")

	return DBPool, nil
}
