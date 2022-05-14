package handlers

func (s *CoreServer) mapErrorToHttpCode(err error) int {
	if err.Error() == "400 Bad Request" {
		return 400
	}
	s.Logger.Warn("Unmapped error: ", err)
	return 500
}

func (s *CoreServer) mapGRPCErrorToHttpCode(err error) int {
	s.Logger.Warn("Unmapped error: ", err)
	return 500
}
