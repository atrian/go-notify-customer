package interfaces

// Logger интерфейс обертка над логгером приложения (Zap)
type Logger interface {
	Fatal(message string, err error)
	Panic(message string, err error)
	Error(message string, err error)
	Warning(message string)
	Info(message string)
	Debug(message ...string)
	Sync()
}
