package net

func syncPublickey(s *Server)*serverEndpoint{
	return &serverEndpoint{
		Methods: []string{ "POST"},
		Handler: syncPublickeyHandler,
		Server:  s,
	}
}

func syncPublickeyHandler(ctx *serverRequestContextImpl) (interface{}, error) {

	return nil, nil
}
