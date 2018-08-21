package sdcustom

// Printer print reporting error
type Printer interface {
	// for inforamation log
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

// Logger print reporting error
type Logger interface {
	Printer

	// for debugging log
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	// for error log
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

type loggerWrapper struct {
	pF  func(v ...interface{})
	pfF func(format string, v ...interface{})
}

// WrapLogger use Printer as Logger
func WrapLogger(v Printer) Logger {
	return &loggerWrapper{
		pF:  v.Print,
		pfF: v.Printf,
	}
}

func (w *loggerWrapper) Debug(v ...interface{}) {
	w.pF(v...)
}
func (w *loggerWrapper) Debugf(format string, v ...interface{}) {
	w.pfF(format, v...)
}
func (w *loggerWrapper) Print(v ...interface{}) {
	w.pF(v...)
}
func (w *loggerWrapper) Printf(format string, v ...interface{}) {
	w.pfF(format, v...)
}
func (w *loggerWrapper) Error(v ...interface{}) {
	w.pF(v...)
}
func (w *loggerWrapper) Errorf(format string, v ...interface{}) {
	w.pfF(format, v...)
}
