package quicklog

var logger = New("./logs/", Config{
	Level:       LEVEL_DEBUG,
	TraceSkip:   2,
	ArchiveLogs: true,
})

func SetDir(logDir string) {
	logger.SetDir(logDir)
}

func GetLevel() LogLevel {
	return logger.GetLevel()
}

func SetLevel(level LogLevel) {
	logger.SetLevel(level)
}

func DoesLogArchives() bool {
	return logger.DoesLogArchives()
}

func SetArchiveLogs(archiveLogs bool) {
	logger.SetArchiveLogs(archiveLogs)
}

func Debugf(format string, a ...any) {
	logger.Debugf(format, a...)
}

func Tracef(format string, a ...any) {
	logger.Tracef(format, a...)
}

func Infof(format string, a ...any) {
	logger.Infof(format, a...)
}

func Warnf(format string, a ...any) {
	logger.Warnf(format, a...)
}

func Errorf(format string, a ...any) {
	logger.Errorf(format, a...)
}

func Fatalf(format string, a ...any) {
	logger.Fatalf(format, a...)
}

func Close() {
	logger.Close()
}
