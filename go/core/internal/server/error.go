package server

import "net/http"

func (s *CoreServer) mapErrorToHttpCode(err error) int {
	switch err.Error() {
	case "400 Bad Request":
		return http.StatusBadRequest
	case "404 Not Found":
		return http.StatusNotFound
	default:
		s.Logger.Warn("Unmapped error: ", err)
		return http.StatusInternalServerError
	}

}
