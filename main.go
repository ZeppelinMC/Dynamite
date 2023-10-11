package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dynamitemc/dynamite/core_commands"
	"github.com/dynamitemc/dynamite/logger"
	"github.com/dynamitemc/dynamite/server"
	"github.com/dynamitemc/dynamite/util"
	"github.com/dynamitemc/dynamite/web"
)

var log = logger.New()
var startTime = time.Now().Unix()

func start(cfg *server.Config) {
	srv, err := server.Listen(cfg, cfg.ServerIP+":"+strconv.Itoa(cfg.ServerPort), log, core_commands.Commands)
	log.Info("Opened TCP server on %s:%d", cfg.ServerIP, cfg.ServerPort)
	if err != nil {
		log.Error("Failed to open TCP server: %s", err)
		os.Exit(1)
	}
	log.Info("Done! (%ds)", time.Now().Unix()-startTime)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		srv.Close()
		os.Exit(0)
	}()

	go srv.ScanConsole()
	err = srv.Start()
	if err != nil {
		log.Error("Failed to start server: %s", err)
		os.Exit(1)
	}
}

func main() {
	log.Info("Starting Dynamite 1.20.1 Server")

	var cfg server.Config
	if err := server.LoadConfig("config.toml", &cfg); err != nil {
		log.Info("%v loading config.toml. Using default config", err)
		cfg = server.DefaultConfig
	}
	log.Debug("Loaded config")

	if !cfg.Online && !util.HasArg("-no_offline_warn") {
		log.Warn("Offline mode is insecure and you should not use it unless for a private server.\nRead https://github.com/DynamiteMC/Dynamite/wiki/Why-you-shouldn't-use-offline-mode")
	}

	if cfg.Web.Enable {
		if !util.HasArg("-nogui") {
			go web.LaunchWebPanel(fmt.Sprintf("%s:%d", cfg.Web.ServerIP, cfg.Web.ServerPort), cfg.Web.Password, log)
		} else {
			log.Warn("Remove the -nogui argument to load the web panel")
		}
	}
	start(&cfg)
}
