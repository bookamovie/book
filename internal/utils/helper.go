package utils

import bookamovierpc "github.com/xoticdsign/bookamovie-proto/gen/go/bookamovie/v2"

func ValidateBookRequest(req *bookamovierpc.BookRequest) bool {
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
