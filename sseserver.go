package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	"strconv"
	"strings"

	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	msgGeneral = "message"
	msgPing    = "ping"
	appName    = "sseserver"
)

var (
	dataStream       = make(chan string)
	inputFile        string
	inputFileContent string
	enableSyslog     bool
	port             int
	host             string
	context          string
	enableDebug      bool
	pingInterval     int
)

func main() {
	flag.StringVar(&inputFile, "input-file", "", "JSON input file *Required*")
	flag.StringVar(&host, "host", "127.0.0.1", "Default binding ip")
	flag.IntVar(&port, "port", 3000, "Port on which events will be sent")
	flag.StringVar(&context, "context", "/", "Context path on which it will listen")
	flag.BoolVar(&enableDebug, "enable-debug", false, "Enable debug mode")
	flag.BoolVar(&enableSyslog, "enable-syslog", false, "Log to syslog")
	flag.IntVar(&pingInterval, "ping-interval", 20, "Interval (in seconds) between ping messages sent to client")
	flag.Parse()

	if inputFile == "" {
		flag.Usage()
		log.Fatal("\nERROR: Input file cannot be empty\n")
	}

	if enableSyslog {
		syslogWriter, err := syslog.New(syslog.LOG_NOTICE, appName)
		defer func() {
			swErr := syslogWriter.Close()
			if swErr != nil {
				log.Fatal(swErr)
			}
		}()
		if err != nil {
			log.Fatalf("Error initializing Syslog Writer: %s\n", err)
		}
		log.SetOutput(syslogWriter)
	}

	readData()
	watchSighupSignal()
	initServer()
}

func initServer() {
	if enableDebug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	endpointRoute(r.Group("/"))
	log.Printf("Using context path: %v \n ", context)

	log.Printf("Listening on %v:%v\n ", host, strconv.Itoa(port))
	err := r.Run(host + ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

}

func endpointRoute(r *gin.RouterGroup) {
	r.GET(context, func(c *gin.Context) {
		initialDataStream := make(chan string)
		go func() { initialDataStream <- inputFileContent }()

		pingStream := make(chan string)
		go periodicPings(pingStream)

		contentType := c.GetHeader("Accept")
		logDebug("Client %s accepts content-type %s\n", c.ClientIP(), contentType)

		if contentType == "application/json" {
			c.Writer.Header().Set("Content-Type", "application/json")

			log.Printf("Requested plain HTTP by %s\n", c.ClientIP())
			c.String(200, inputFileContent)
		} else if contentType == "text/event-stream" || contentType == "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Keep-Alive,X-Requested-With,Cache-Control,Content-Type,Last-Event-ID")

			log.Printf("Requested SSE stream by %s\n", c.ClientIP())
			c.Stream(func(w io.Writer) bool {
				select {
				case initialData := <-initialDataStream:
					logDebug("Sending initial data message to %s\n", c.ClientIP())
					c.SSEvent(msgGeneral, initialData)
				case data := <-dataStream:
					logDebug("Sending new data message to %s\n", c.ClientIP())
					c.SSEvent(msgGeneral, data)
				case ping := <-pingStream:
					logDebug("Sending ping message to %s\n", c.ClientIP())
					c.SSEvent(msgPing, ping)
				}
				return true
			})
		} else {
			log.Printf("Unsupported content-type requested by %s\n", c.ClientIP())
			c.String(http.StatusBadRequest, "Unsupported content-type %s", contentType)
		}
	})
}

func logDebug(formattedMessage string, v ...interface{}) {
	if enableDebug {
		log.Printf(formattedMessage, v)
	}
}

func periodicPings(pingStream chan string) {
	for {
		pingStream <- msgPing
		time.Sleep(time.Duration(pingInterval) * time.Second)
	}
}

func watchSighupSignal() {
	sighupChannel := make(chan os.Signal, 1)
	signal.Notify(sighupChannel, syscall.SIGHUP)

	go func() {
		for sig := range sighupChannel {
			log.Printf("%v signal, sending file %v into stream\n", sig, inputFile)
			dataStream <- readData()
		}
	}()
}

func readData() string {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatal("Can't read inputfile: ", err)
	}

	var jsonObj interface{}
	err = json.Unmarshal(data, &jsonObj)
	if err != nil {
		log.Fatal("Inputfile isn't JSON: ", err)
	}

	inputFileContent = strings.Replace(string(data), "\n", "", -1) // return the content of the file and strip newlines
	return inputFileContent
}
