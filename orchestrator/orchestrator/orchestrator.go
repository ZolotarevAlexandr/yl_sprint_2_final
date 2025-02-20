package orchestrator

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func RunOrchestrator() {
	mux := http.NewServeMux()

	mux.Handle("/api/v1/calculate", LoggingMiddleware(http.HandlerFunc(handleCalculate)))
	mux.Handle("/api/v1/expressions", LoggingMiddleware(http.HandlerFunc(expressionsHandler)))
	mux.Handle("/internal/task", LoggingMiddleware(http.HandlerFunc(internalTaskHandler)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Orchestrator is running on %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
