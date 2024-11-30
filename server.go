package main
import (
	"fmt"
	"boards/internals/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"boards/db"
)
var Client *mongo.Client
func main(){
	app := fiber.New()
	app.Use(cors.New())
	var err error
    _, err = db.DBConnection() // No need to assign the client to a local variable

    // If an error occurs in DB connection
    if err != nil {
        log.Fatalf("Error establishing DB connection: %v", err)
    }

    fmt.Println("Client is connected:")
	app.Get("/boards", handlers.GetBoardDetailsHandler)
	app.Get("/boardTasks", handlers.GetBoardTasksHandler)
	port := app.Listen(":7000")
	if port != nil {
		fmt.Println("Error in connecting the Port")
	}
}