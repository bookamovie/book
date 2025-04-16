package utils

import bookrpc "github.com/bookamovie/proto/gen/go/book/v3"

// ValidateBookRequest() validates the fields in the BookRequest to ensure all required information is provided.
//
// It checks whether the cinema name, movie title, session seat, screen, and session date are properly set.
func ValidateBookRequest(req *bookrpc.BookRequest) bool {
	switch {
	case req.GetCinema().GetName() == "":
		return false

	case req.GetCinema().GetLocation() == "":
		return false

	case req.GetMovie().GetTitle() == "":
		return false

	case req.GetSession().GetSeat() == 0:
		return false

	case req.GetSession().GetScreen() == 0:
		return false

	case req.GetSession().GetDate().AsTime().IsZero():
		return false
	}

	return true
}
