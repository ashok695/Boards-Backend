// dbconnection.go
package db

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
    "time"
)

// Export the client variable so it can be accessed in other packages
var Client *mongo.Client

// DBConnection establishes a connection to the MongoDB database and returns the client
func DBConnection() (*mongo.Client, error) {
    fmt.Println("Initializing DB connection...")

    MONGOURL := "mongodb+srv://qas-app-rwuser:fSXQNZ7G9B38n4xi@quality-app.ss86w.mongodb.net/"
    options := options.Client().ApplyURI(MONGOURL)

    // Create a new context with timeout for DB connection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    Client, err = mongo.Connect(ctx, options)
    if err != nil {
        log.Println("Error connecting to DB:", err)
        return nil, err
    }

    // Ping the database to ensure connection is successful
    pingErr := Client.Ping(ctx, nil)
    if pingErr != nil {
        log.Println("Error pinging the database:", pingErr)
        return nil, pingErr
    }

    fmt.Println("DB CONNECTION SUCCESSFUL")
    return Client, nil
}
