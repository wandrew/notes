package server

import (
	"encoding/json"
	"github.com/wandrew/mhmmm/notes/api"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
)

type Server struct {
	nc api.NoteServiceClient
}

func NewServer(c *grpc.ClientConn) *Server{
	nsc := api.NewNoteServiceClient(c)
	return &Server{
		nc: nsc,
	}
}

// CreateNote calls gRPC Api for creating a note
func (s *Server) CreateNote(w http.ResponseWriter, r *http.Request) {
	apiRequest := api.CreateNoteRequest{}

	httpRequestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		JSONError(w, err, 500)
	}
	json.Unmarshal(httpRequestBody, apiRequest)

	resp, err := s.nc.Create(r.Context(),&apiRequest)

	json.NewEncoder(w).Encode(resp)
}

// ListNotes calls gRPC Api for searching with no params.
func (s *Server) ListNotes(w http.ResponseWriter, r *http.Request) {
	resp, err := s.nc.Search(r.Context(), &api.SearchNotesRequest{})
	if err != nil {
		JSONError(w, err, 500)
	}
	json.NewEncoder(w).Encode(resp)
}

// JSONError is a helper function for emitting http errors when this server runs into a problem
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}