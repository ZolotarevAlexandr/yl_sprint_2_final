package orchestrator

import (
	"encoding/json"
	"net/http"
	"strings"
)

// handlePing handles GET /api/v1/ping healthcheck endpoint.
// It should always return 200 with message "orchestrator is up and running"
func handlePing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "orchestrator is up and running"})
}

// handleCalculate processes POST /api/v1/calculate to add a new expression.
func handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, "invalid data", http.StatusUnprocessableEntity)
		return
	}
	expr, err := BuildExpressionTasks(req.Expression)
	if err != nil {
		http.Error(w, "error processing expression", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": expr.ID})
}

// internalTaskHandler handles agent requests: GET for retrieving a task and POST for submitting the result.
func internalTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleGetTask(w, r)
	} else if r.Method == http.MethodPost {
		handlePostTask(w, r)
	}
}

// handleListExpressions returns a list of all expressions.
func handleListExpressions(w http.ResponseWriter, r *http.Request) {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	var exprList []Expression
	for _, expr := range expressionsStore {
		exprList = append(exprList, *expr)
	}
	json.NewEncoder(w).Encode(map[string]any{"expressions": exprList})
}

// handleGetExpression returns a specific expression by its id.
func handleGetExpression(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	storeMutex.Lock()
	expr, ok := expressionsStore[id]
	storeMutex.Unlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"expression": expr})
}

// handleGetTask returns a task to the agent for computation.
func handleGetTask(w http.ResponseWriter, r *http.Request) {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	for _, task := range tasksStore {
		if task.Status == "pending" && updateTaskDependencies(task) {
			task.Status = "running"
			resp := map[string]any{
				"task": map[string]any{
					"id":             task.ID,
					"arg1":           *task.Arg1,
					"arg2":           *task.Arg2,
					"operation":      task.Operator,
					"operation_time": task.OperationTime,
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	http.Error(w, "no task", http.StatusNotFound)
}

// handlePostTask accepts the result from the agent and updates the task status.
func handlePostTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" {
		http.Error(w, "invalid data", http.StatusUnprocessableEntity)
		return
	}
	storeMutex.Lock()
	task, ok := tasksStore[req.ID]
	if !ok {
		storeMutex.Unlock()
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	if task.Status != "running" {
		storeMutex.Unlock()
		http.Error(w, "task not in running state", http.StatusUnprocessableEntity)
		return
	}
	task.Status = "done"
	task.Result = &req.Result
	// If this is the root task, update the expression status
	expr, exists := expressionsStore[task.ExpressionID]
	if exists && expr.RootTaskID == task.ID {
		expr.Status = "done"
		expr.Result = &req.Result
	}
	storeMutex.Unlock()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "result recorded"})
}
