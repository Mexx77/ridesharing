package server

import (
	"context"
	"fmt"
	"github.com/Mexx77/ridesharing/config"
	"github.com/Mexx77/ridesharing/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

type server struct {
	users		*mongo.Collection
	rides		*mongo.Collection
	mongoClient *mongo.Client
	config      *config.Config
}

func NewServer() {
	conf := config.GetConfig()

	// MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf("mongodb+srv://%s:%s@%s", conf.MongoUsername, conf.MongoPw, conf.MongoHost),
	))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	logging.Info.Println("Connected to MongoDB!")

	server := &server{
		rides: client.Database("ridesharing").Collection("rides"),
		users: client.Database("ridesharing").Collection("users"),
		mongoClient: client,
		config: conf,
	}

	defer server.closeMongoConnection()
	server.startHttpServer()
}

func (s *server) startHttpServer() {
	httpServer := &http.Server{
		Addr:     fmt.Sprintf(":%s", s.config.Port),
		ErrorLog: logging.Error,
	}

	logging.Info.Print("Starting http server on port ", s.config.Port)
	s.routes()
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}

}

func (s *server) addCORSHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.config.Environment == config.DevEnvironment {
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, Authorization")
			w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		h(w, r)
	}
}

func (s *server) closeMongoConnection() {
	err := s.mongoClient.Disconnect(context.TODO())
	if err != nil {
		logging.Error.Println(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}