package utils

import bookrpc "github.com/xoticdsign/bookamovie-proto/gen/go/book/v3"

func ValidateBookRequest(req *bookrpc.BookRequest) bool {
	switch {
	case req.GetCinema().GetName() == "":
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
