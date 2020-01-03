package add

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporter/trace/stdout"
	"go.opentelemetry.io/otel/plugin/othttp"
)

func Run() {
	std, err := stdout.NewExporter(stdout.Options{PrettyPrint: true})
	if err != nil {
		log.Fatal(err)
	}

	traceProvider, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(std))
	if err != nil {
		log.Fatal(err)
	}

	global.SetTraceProvider(traceProvider)

	mux := http.NewServeMux()
	mux.Handle("/", othttp.NewHandler(http.HandlerFunc(addHandler), "add", othttp.WithPublicEndpoint()))

	log.Println("Initializing server...")
	err = http.ListenAndServe(":3001", mux)
	if err != nil {
		log.Fatalf("Could not initialize server: %s", err)
	}
}

func addHandler(w http.ResponseWriter, req *http.Request) {
	var x int
	var y int

	var err error

	if (len(req.URL.Query()["x"]) != 1) && (len(req.URL.Query()["y"]) != 1) {
		err = fmt.Errorf("Invalid arguments")
	} else {
		x, err = strconv.Atoi(req.URL.Query()["x"][0])
		y, err = strconv.Atoi(req.URL.Query()["y"][0])
	}
	if err != nil {
		fmt.Fprintf(w, "Failed to parse querystring")
		w.WriteHeader(503)
		return
	}
	fmt.Fprintf(w, "%d", x+y)
}
