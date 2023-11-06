package app

type App interface {
	Serve() error
	Shutdown()
}
