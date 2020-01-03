package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporter/trace/stdout"
	"go.opentelemetry.io/otel/plugin/httptrace"
	"go.opentelemetry.io/otel/plugin/othttp"
)

var services Config

func Start() {
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
	mux.Handle("/", othttp.NewHandler(http.HandlerFunc(rootHandler), "root", othttp.WithPublicEndpoint()))
	mux.Handle("/calculate", othttp.NewHandler(http.HandlerFunc(calcHandler), "calculate", othttp.WithPublicEndpoint()))
	services = GetServices()

	log.Println("Initializing server...")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("Could not initialize server: %s", err)
	}
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	trace.CurrentSpan(ctx).AddEvent(ctx, "called root handler, getting discovered services")
	fmt.Fprintf(w, "%s", services)
}

func calcHandler(w http.ResponseWriter, req *http.Request) {
	var res int

	calcRequest, err := ParseCalcRequest(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch method := calcRequest.Method; strings.ToLower(method) {
	case "add":
		res, err = callAdd(req.Context(), calcRequest.Operands[0], calcRequest.Operands[1])
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s", res)
}

func callAdd(ctx context.Context, x int, y int) (int, error) {
	client := http.DefaultClient
	url := fmt.Sprintf("http://localhost:3001/add?x=%d&y=%d", x, y)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	ctx, req = httptrace.W3C(ctx, req)
	httptrace.Inject(ctx, req)
	res, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return -1, err
	}
	resp, err := strconv.Atoi(string(body))
	if err != nil {
		return -1, err
	}
	return resp, nil
}
