package quicklog

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type openedFile struct {
	file    *os.File
	created time.Time
}

func (file *openedFile) File() *os.File {
	return file.file
}

func (file *openedFile) Created() time.Time {
	return file.created
}

func (file *openedFile) IsOK() bool {
	return file.file != nil
}

func (file *openedFile) IsToday() bool {
	now := time.Now()
	created := file.created
	return created.YearDay() == now.YearDay() && created.Year() == now.Year()
}

func (file *openedFile) Close() {
	file.file.Close()
	file.file = nil
}

const timestampFormat = "2006-01-02; 15:04:05"

type Config struct {
	// If the len(Title) > 0, the title will be prepended before message content
	Label string

	// Lowest log level that can be handled
	Level LogLevel

	// The number of callstack frames to skip. This would be the argument passed to runtime.Caller(skip int)
	// By default, this is 1 such that the Trace will log where the statement was executed
	TraceSkip int

	// Enables/disabled writing the current log to a file. The contents of the log file is overwritten
	// during the first Logger.Write(). If ArchiveLogs is true, this option is ignored
	WriteLogFile bool

	// Enables/disables log archiving
	ArchiveLogs bool

	// If File is nil, it will default to os.Stdout
	File *os.File
}

type Logger struct {
	io.Writer

	dir         string
	logFilePath string
	logFile     openedFile

	format string

	config Config
}

func New(logDir string, config ...Config) *Logger {
	logger := new(Logger)
	logger.SetDir(logDir)
	logger.config.TraceSkip = 1
	logger.config.ArchiveLogs = true // by default, create log archives

	if len(config) > 0 {
		logger.config = config[0]

		if logger.config.File == nil {
			logger.config.File = os.Stdout
		}
	}

	logger.updateFormat()
	return logger
}

func (logger *Logger) updateFormat() {
	if len(logger.config.Label) > 0 {
		logger.format = "[%[1]s; %[2]s] {%[3]s} %[4]s\n"
	} else {
		logger.format = "[%[1]s; %[2]s] %[4]s\n"
	}
}

func (logger *Logger) Close() {
	err := logger.ArchiveCurrentLog()
	if err != nil {
		panic(err)
	}
	logger.logFile.Close()
}

func (logger *Logger) openLogFile() openedFile {
	err := os.MkdirAll(logger.dir, 0775)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(logger.logFilePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0775)
	if err != nil {
		panic(err)
	}

	stat, err := os.Stat(logger.logFilePath)
	if err != nil {
		panic(err)
	}

	return openedFile{
		file:    file,
		created: stat.ModTime(),
	}
}

func (logger *Logger) SetDir(logDir string) {
	logger.dir = logDir
	logger.logFilePath = fmt.Sprint(logDir, "current.log")
}

func (logger *Logger) GetLevel() LogLevel {
	return logger.config.Level
}

func (logger *Logger) SetLevel(level LogLevel) {
	logger.config.Level = level
}

func (logger *Logger) DoesLogArchives() bool {
	return logger.config.ArchiveLogs
}

func (logger *Logger) ArchiveLogs(archiveLogs bool) {
	logger.config.ArchiveLogs = archiveLogs
}

func (logger *Logger) WriteLogFile(writeLog bool) {
	logger.config.WriteLogFile = writeLog
}

func (logger *Logger) File() *os.File {
	if logger.config.File == nil {
		return os.Stdout
	}
	return logger.config.File
}

func (logger *Logger) ArchiveCurrentLog() error {
	if !logger.config.ArchiveLogs {
		return nil
	}

	year, month, day := logger.logFile.Created().Date()
	filename := fmt.Sprintf("%s%d-%02d-%02d", logger.dir, year, month, day)

	files, err := filepath.Glob(fmt.Sprint(filename, "*.log.gz"))
	if err != nil {
		return err
	}

	archive, err := os.OpenFile(fmt.Sprintf("%s_%d.log.gz", filename, len(files)+1), os.O_WRONLY|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer archive.Close()

	writer := gzip.NewWriter(archive)
	defer writer.Close()

	logger.logFile.File().Seek(0, io.SeekStart)
	bytes, err := io.ReadAll(logger.logFile.File())
	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)

	return err
}

func join(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (logger *Logger) Write(p []byte) (n int, err error) {
	var aErr, bErr error
	n, aErr = logger.File().Write(p)

	if !logger.config.WriteLogFile && !logger.config.ArchiveLogs {
		return n, aErr
	}

	if !logger.logFile.IsOK() {
		logger.logFile = logger.openLogFile()
	}

	if logger.logFile.IsToday() {
		_, bErr = logger.logFile.File().Write(p)
	} else {
		logger.ArchiveCurrentLog()
		logger.logFile.Close()
		logger.logFile = logger.openLogFile()
	}

	return n, join(aErr, bErr)
}

func (logger *Logger) WriteLog(message string, level LogLevel) {
	timestamp := time.Now().Format(timestampFormat)
	fmt.Fprintf(logger, logger.format, LOG_PREFIX[level], timestamp, logger.config.Label, message)
}

func (logger *Logger) Logf(level LogLevel, format string, a ...any) {
	if level < logger.config.Level {
		return
	}

	logger.WriteLog(fmt.Sprintf(format, a...), level)
}

func (logger *Logger) Debugf(format string, a ...any) {
	logger.Logf(LEVEL_DEBUG, format, a...)
}

func (logger *Logger) Tracef(format string, a ...any) {
	_, filename, line, _ := runtime.Caller(logger.config.TraceSkip)
	traceMessage := fmt.Sprintf("[%s:%d]", filename, line)
	logger.Logf(LEVEL_TRACE, fmt.Sprint(traceMessage, " ", format), a...)
}

func (logger *Logger) Infof(format string, a ...any) {
	logger.Logf(LEVEL_INFO, format, a...)
}

func (logger *Logger) Warnf(format string, a ...any) {
	logger.Logf(LEVEL_WARN, format, a...)
}

func (logger *Logger) Errorf(format string, a ...any) {
	logger.Logf(LEVEL_ERROR, format, a...)
}

func (logger *Logger) Error(err error) {
	logger.Logf(LEVEL_ERROR, err.Error())
}

func (logger *Logger) Fatalf(format string, a ...any) {
	logger.Logf(LEVEL_FATAL, format, a...)
	panic(fmt.Sprintf(format, a...))
}

func (logger *Logger) FatalErr(err error) {
	logger.Logf(LEVEL_FATAL, err.Error())
	panic(err)
}
