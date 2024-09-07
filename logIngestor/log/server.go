package log

// Server interface
type Server interface {
	Start() error
	Close() error
}
