package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"os"

	_ "github.com/lib/pq"
	"github.com/mundipagg/boleto-api/app"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/robot"
)

var (
	env          string
	processID    = os.Getpid()
	totalProcs   = runtime.NumCPU()
	devMode      = flag.Bool("dev", false, "-dev To run in dev mode")
	mockMode     = flag.Bool("mock", false, "-mock To run mock requests")
	disableLog   = flag.Bool("nolog", false, "-nolog disable seq log")
	airPlaneMode = flag.Bool("airplane-mode", false, "-airplane-mode run api in dev, mock and nolog mode")
	mockOnly     = flag.Bool("mockonly", false, "-mockonly run just mock service")
)

func init() {
	createPIDfile()
	// Configure signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go handleSignal(c)
}

func handleSignal(c chan os.Signal) {
	<-c
	config.Stop()
	log.Info("Quiting BoletoApi")
	log.Close()
	//db.GetDB().Close()
	fmt.Println("Done")
	os.Exit(1)
}

func createPIDfile() {
	p := strconv.Itoa(processID)
	ioutil.WriteFile("boletoapi.pid", []byte(p), 0644)
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if *mockOnly {
		w := make(chan int)
		config.Install(true, true, true)
		robot.GoRobots()
		<-w
	} else {
		params := app.NewParams()
		if *airPlaneMode {
			params.DevMode = true
			params.DisableLog = true
			params.MockMode = true
			env = strconv.FormatBool(params.DevMode)
		} else {
			params.DevMode = *devMode
			params.DisableLog = *disableLog
			params.MockMode = *mockMode
			env = strconv.FormatBool(params.DevMode)
		}
		logo(env)
		app.Run(params)
	}

}

func logo(env string) {
	l := `
$$$$$$$\            $$\            $$\                $$$$$$\            $$\ 
$$  __$$\           $$ |           $$ |              $$  __$$\           \__|
$$ |  $$ | $$$$$$\  $$ | $$$$$$\ $$$$$$\    $$$$$$\  $$ /  $$ | $$$$$$\  $$\ 
$$$$$$$\ |$$  __$$\ $$ |$$  __$$\\_$$  _|  $$  __$$\ $$$$$$$$ |$$  __$$\ $$ |
$$  __$$\ $$ /  $$ |$$ |$$$$$$$$ | $$ |    $$ /  $$ |$$  __$$ |$$ /  $$ |$$ |
$$ |  $$ |$$ |  $$ |$$ |$$   ____| $$ |$$\ $$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |
$$$$$$$  |\$$$$$$  |$$ |\$$$$$$$\  \$$$$  |\$$$$$$  |$$ |  $$ |$$$$$$$  |$$ |
\_______/  \______/ \__| \_______|  \____/  \______/ \__|  \__|$$  ____/ \__|
                                                               $$ |          
                                                               $$ |          
                                                               \__|          
	`
	fmt.Println(l)
	fmt.Println("Version: " + config.Get().Version)
	fmt.Println("DevMode: " + env)
}
