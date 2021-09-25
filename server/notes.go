package server

import (
	"context"
	"log"
	"time"

	"github.com/wandrew/mhmmm/notes/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
)

type NoteService struct {
	mongo      *mongo.Client
	collection *mongo.Collection
}

func NewNoteService(mongo *mongo.Client) *NoteService {
	return &NoteService{
		mongo:      mongo,
		collection: mongo.Database("mhmmm").Collection("notes"),
	}
}

// Create implements the api.NoteServiceClient interface, a protobuf api creating the endpoint to create new notes
func (s *NoteService) Create(ctx context.Context, request *api.CreateNoteRequest) (*api.CreateNoteResponse, error) {
	uuid := uuid.New()
	// This is simple, and simple is good.
	//  For more complex operations, I might consider creating a `data` pkg in this service.  such a package could
	//    include data-models and logic for each database/collection.  For this though, I wouldn't complicate the
	//    application.  I like to make it exactly as complicated as needed.
	_, err := s.collection.InsertOne(ctx, bson.D{
		{"id", uuid.String()},
		{"title", request.Note.Title},
		{"content", request.Note.Content},
		{"created", time.Now().String()},
	})

	if err != nil {
		return nil, err
	}

	return &api.CreateNoteResponse{Id: uuid.String()}, nil
}

// Search implements the api.NoteServiceClient interface, a protobuf api interface, and creates the endpoint to search for notes.
func (s *NoteService) Search(ctx context.Context, request *api.SearchNotesRequest) (*api.SearchNotesResponse, error) {
	c, err := s.collection.Find(ctx, bson.D{})
	defer c.Close(ctx)
	if err != nil {
		return nil, err
	}

	resp := api.SearchNotesResponse{}
	for c.Next(ctx) {
		note := api.Note{}
		err := c.Decode(&note)
		if err != nil {
			log.Printf("could not decode note from mongo: %v", err)
			continue
		}
		resp.Notes = append(resp.Notes, &note)
	}

	return &resp, nil
}
