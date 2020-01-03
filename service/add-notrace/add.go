package addnotrace

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Run() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(addHandler))

	log.Println("Initializing server...")
	err := http.ListenAndServe(":3001", mux)
	if err != nil {
		log.Fatalf("Could not initialize server: %s", err)
	}
}

func addHandler(w http.ResponseWriter, req *http.Request) {
	values := strings.Split(req.URL.Query()["o"][0], ",")
	var res int
	for _, n := range values {
		i, err := strconv.Atoi(n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res += i
	}
	fmt.Fprintf(w, "%d", res)
}
