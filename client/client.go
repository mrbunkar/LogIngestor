package client

type Logger struct {
	Addr   string
	Format interface{}
}

func NewLogger(addr string, format interface{}) (*Logger, error) {
	logger := &Logger{
		Addr:   addr,
		Format: format,
	}

	if err := logger.HandShake(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *Logger) HandShake() error {
	return nil
}

func (l *Logger) Info() {}

func (l *Logger) Error() {}
