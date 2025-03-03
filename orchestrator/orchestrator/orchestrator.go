package orchestrator

import (
	"fmt"
	"log"
	"net/http"
)

func RunOrchestrator() {
	mux := http.NewServeMux()

	mux.Handle("/api/v1/ping", ErrorHandlingMiddleware(LoggingMiddleware(http.HandlerFunc(handlePing))))
	mux.Handle("/api/v1/calculate", ErrorHandlingMiddleware(LoggingMiddleware(http.HandlerFunc(handleCalculate))))
	mux.Handle("/api/v1/expressions", ErrorHandlingMiddleware(LoggingMiddleware(http.HandlerFunc(handleListExpressions))))
	mux.Handle("/api/v1/expressions/", ErrorHandlingMiddleware(LoggingMiddleware(http.HandlerFunc(handleGetExpression))))
	mux.Handle("/internal/task", ErrorHandlingMiddleware(LoggingMiddleware(http.HandlerFunc(internalTaskHandler))))

	fmt.Printf("Orchestrator is running on %s\n", Port)
	if err := http.ListenAndServe(":"+Port, mux); err != nil {
		log.Fatal(err)
	}
}
