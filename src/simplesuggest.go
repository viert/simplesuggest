package main

import (
	"config"
	"flag"
	logging "github.com/op/go-logging"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"web"
)

var (
	logfileName string
	logfile     *os.File
	log         *logging.Logger = logging.MustGetLogger("suggest")
)

const (
	DEFAULT_CONFIG_FILE = "/etc/simplesuggest.conf"
)

func reopenLog() {
	var err error
	if logfile != nil {
		logfile.Close()
	}
	logfile, err = os.OpenFile(logfileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return
	}
	backend := logging.NewLogBackend(logfile, "", 0)
	logging.SetBackend(backend)
}

func hupCatcher() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP)
	for _ = range c {
		log.Debug("HUP signal catched, reopening logfile %s", logfileName)
		reopenLog()
	}
}

func main() {
	var confFile, format string
	flag.StringVar(&confFile, "c", DEFAULT_CONFIG_FILE, "suggest daemon config filename")
	flag.Parse()

	conf := config.Load(confFile)
	switch conf.Log {
	case "syslog":
		backend, err := logging.NewSyslogBackend("suggest")
		if err != nil {
			log.Error(err.Error())
			return
		}
		logging.SetBackend(backend)
		format = "%{program} %{level} %{message}"
	case "":
		backend := logging.NewLogBackend(os.Stdout, "", 0)
		backend.Color = true
		logging.SetBackend(backend)
		format = "%{color:reset}[%{time:2006-01-02 15:04:05}] %{color}%{level} %{color:reset}%{message}"
	default:
		logfileName = conf.Log
		reopenLog()
		go hupCatcher()
		format = "[%{time:2006-01-02 15:04:05}] %{level} %{message}"
	}

	logging.SetFormatter(logging.MustStringFormatter(format))
	logging.SetLevel(conf.LogLevel, "suggest")

	log.Debugf("Configuring runtime: GCPercent(%d), MaxCores(%d), MaxThreads(%d)", conf.GCPercent, conf.MaxCores, conf.MaxThreads)
	runtime.GOMAXPROCS(conf.MaxCores)
	debug.SetGCPercent(conf.GCPercent)
	debug.SetMaxThreads(conf.MaxThreads)

	server := web.NewServer()
	server.Start(conf.Host, conf.Port)

}
