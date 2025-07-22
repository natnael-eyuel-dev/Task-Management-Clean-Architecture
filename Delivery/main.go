package main

// imports
import (
	"context";
	"log";
	"time";
	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Delivery/routers";
	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Infrastructure";
	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Repositories";
	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Usecases";
	"go.mongodb.org/mongo-driver/mongo";
	"go.mongodb.org/mongo-driver/mongo/options";
)

// entry point of the Task Management application
func main() {

	// setup mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)       // set timeout
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)       // disconnect

	db := client.Database("taskmanager")
	taskCol := db.Collection("tasks")         // initialize task collection
	userCol := db.Collection("users")         // initialize user collection

	jwtservice, _ := infrastructure.NewJWTService()              // setup jwt service infrastructure
	passwordService := infrastructure.NewPasswordService()       // setup password service infrastructure

	taskRepo := repositories.NewTaskRepository(taskCol)          // setup task repositorie
	userRepo := repositories.NewUserRepository(userCol)          // setup user repositorie

	taskUC := usecases.NewTaskUseCase(taskRepo)                                    // setup task use case
	userUC := usecases.NewUserUseCase(userRepo, jwtservice, passwordService)       // setup user use case

	router := routers.SetupRouter(taskUC, userUC, jwtservice)       // initialize the router with all configured routes

	// start the server on port 8080
	router.Run(":8080")                        
	log.Println("Starting server on :8080")
}
