package server

type Module interface {
	GetRoutes() []Route
}
