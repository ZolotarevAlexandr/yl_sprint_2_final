package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ZolotarevAlexandr/yl_sprint_2_final/calculator/calculator"
)

// worker is a goroutine that continuously requests tasks.
func worker(workerID int) {
	orchestratorPort := os.Getenv("ORCHESTRATOR_PORT")
	if orchestratorPort == "" {
		orchestratorPort = "8080"
	}
	client := &http.Client{}
	for {
		resp, err := client.Get(fmt.Sprintf("http://localhost:%s/internal/task", orchestratorPort))
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			// No tasks available
			resp.Body.Close()
			time.Sleep(500 * time.Millisecond)
			continue
		}
		var taskResp struct {
			Task struct {
				ID            string  `json:"id"`
				Arg1          float64 `json:"arg1"`
				Arg2          float64 `json:"arg2"`
				Operation     string  `json:"operation"`
				OperationTime int     `json:"operation_time"`
			} `json:"task"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
			resp.Body.Close()
			time.Sleep(500 * time.Millisecond)
			continue
		}
		resp.Body.Close()
		task := taskResp.Task
		// Simulate long computation time
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		// Compute the operation using EvaluateOperation from the calculator package.
		result, err := calculator.EvaluateOperation(task.Operation, task.Arg1, task.Arg2)
		if err != nil {
			log.Printf("Worker %d: error computing task %s: %v", workerID, task.ID, err)
			continue
		}
		// Send the result back to the orchestrator.
		payload, _ := json.Marshal(map[string]any{
			"id":     task.ID,
			"result": result,
		})
		_, err = client.Post(fmt.Sprintf("http://localhost:%s/internal/task", orchestratorPort), "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Printf("Worker %d: error posting result for task %s: %v", workerID, task.ID, err)
			continue
		}
		log.Printf("Worker %d: completed task %s with result %f", workerID, task.ID, result)
	}
}

func RunAgent() {
	cpStr := os.Getenv("COMPUTING_POWER")
	cp, err := strconv.Atoi(cpStr)
	if err != nil || cp < 1 {
		cp = 1
	}
	for i := 0; i < cp; i++ {
		go worker(i)
	}
	fmt.Printf("Started agent with %d workers\n", cp)
	// Infinite wait
	select {}
}
