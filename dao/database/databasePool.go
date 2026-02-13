package database

import "github.com/mt1976/frantic-core/logHandler"

// addConnectionToPool adds a database connection to the connection pool
func addConnectionToPool(db DB, key string) {
	logHandler.Database.Printf("[CON]{CONNECTION}{POOL} Adding [%v] to connection pool (%v)", key, db.databaseName)
	if len(connectionPool) >= connectionPoolMaxSize {
		logHandler.Database.Panicf("[CON]{CONNECTION}{POOL} Connection pool full [%v]", connectionPoolMaxSize)
		return
	}
	connectionPool[key] = &db
	logHandler.Database.Printf("[CON]{CONNECTION}{POOL} Connection pool [size=%v]", len(connectionPool))
}

// releaseFromConnectionPool removes a database connection from the connection pool
func releaseFromConnectionPool(db *DB) {
	logHandler.Database.Printf("[CON]{CONNECTION}{POOL} Removing [%v] from connection pool (%v)", db.Name, db.databaseName)
	for key, value := range connectionPool {
		logHandler.Database.Printf("[CON]{CONNECTION}{POOL}  Connection Pool [%v] [%v] [codec=%v]", key, value.databaseName, value.connection.Node.Codec().Name())
	}
	connectionPool[db.Name] = nil
	delete(connectionPool, db.Name)
	logHandler.Database.Printf("[CON]{CONNECTION}{POOL} Connection pool [size=%v]", len(connectionPool))
	for key, value := range connectionPool {
		logHandler.Database.Printf("[CON]{CONNECTION}{POOL}  Connection Pool [%v] [%v] [codec=%v]", key, value.databaseName, value.connection.Node.Codec().Name())
	}
}
