package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"context"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/wandrew/mhmmm/notes/api"
	"github.com/wandrew/mhmmm/notes/server"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"

)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Ordinarily, I would have a package called config and use some kind of ci pipeline to do an envstringsub operation
	//  to create configuration values based on an environment.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	notesService := server.NewNoteService(client)
	api.RegisterNoteServiceServer(grpcServer, notesService)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("could not start grpc service: %v", err)
	}

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	s := server.NewServer(conn)



	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request){ fmt.Print("a simple notes microservice")})
	router.HandleFunc("/create", s.CreateNote)
	router.HandleFunc("/list", s.ListNotes)

	log.Fatal(http.ListenAndServe(":8080", router))
}
