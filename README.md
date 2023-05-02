# Quick Log

Quick Log is a small logging that package supports logging levels, log archiving, and custom loggers.

## Installing

After initializing your go module, run:

```go get -u github.com/I-Am-Dench/quick-log/v2```

## Log Archiving

This package comes with an automatic archiving of log files built-in. When the logger is created, it will generate a file called `current.log` inside of the configured log directory (`./log/` by default). Whenever the logger is closed, the contents of `current.log` are compressed through gzip and written to the log directory as `yyyy-mm-dd_#.log.gz` where `#` is the number of logs archived for that day. The `#` is added so that if the logger is closed multiple times in a day, each archive does not override one another.

Whenever any logger write function is called, i.e. `Logger.Infof()`, `Logger.Errorf()`, etc., if the current day is not equal to when `current.log` was created, the log will be archived.

Log archiving can be disabled by calling `Logger.SetArchiveLogs(false)`

## Examples

### Global Logger

**Source:** `main.go`

~~~go
package main

import log "github.com/I-Am-Dench/quick-log/v2"

func main() {
    log.Debugf("Debug log")
    log.Tracef("Trace log")
    log.Infof("Info log")
    log.Warnf("Warn log")
    log.Errorf("Error log")
    // log.Fatalf("Fatal log")

    log.Close()
}
~~~

**Output**

~~~
[D; 1970-01-01; 00:00:00] Debug log
[T; 1970-01-01; 00:00:00] [main.go:7] Trace log
[I; 1970-01-01; 00:00:00] Info log
[W; 1970-01-01; 00:00:00] Warn log
[E; 1970-01-01; 00:00:00] Error log
~~~

`Logger.Fatalf` will call `panic()` after logger has finished handling the log message.

### Custom Logger

**Source:** `main.go`

~~~go
package main

import log "github.com/I-Am-Dench/quick-log/v2"

func main() {
    logger1 := log.New("./dir1/logs/")
    defer logger1.Close()

    logger2 := log.New("./dir2/logs/", log.Config{
        Label: "LOGGER2",
    })
    defer logger2.Close()
    logger2.SetLevel(log.LEVEL_INFO)

    logger1.Debugf("Logger1 - Debug")
    logger1.Infof("Logger1 - Info")

    logger2.Debugf("Logger2 - Debug") // Ignored since log.LEVEL_DEBUG < log.LEVEL_INFO
    logger2.Infof("Logger2 - Info")
}
~~~

**Output**

~~~
[D; 1970-01-01; 00:00:00] Logger1 - Debug
[I; 1970-01-01; 00:00:00] Logger1 - Info
[I; 1970-01-01; 00:00:00] {LOGGER 2} Logger2 - Info
~~~

Loggers, by default, have a log level of log.LEVEL_DEBUG.

## Log Levels

The priority of log levels goes as follows:

log.LEVEL_DEBUG < log.LEVEL_TRACE < log.LEVEL_INFO < log.LEVEL_WARN < log.LEVEL_ERROR < log.LEVEL_FATAL