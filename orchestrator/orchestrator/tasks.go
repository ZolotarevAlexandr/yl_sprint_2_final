package orchestrator

import (
	"errors"
	"os"
	"strconv"
	"sync"

	"github.com/ZolotarevAlexandr/yl_sprint_2_final/calculator/calculator"
	"github.com/google/uuid"
)

var (
	expressionsStore = make(map[string]*Expression)
	tasksStore       = make(map[string]*Task)
	storeMutex       sync.Mutex
)

// Expression represents an expression submitted by the user.
type Expression struct {
	ID         string   `json:"id"`
	Expr       string   `json:"expression"`
	Status     string   `json:"status"` // "pending" or "done"
	Result     *float64 `json:"result,omitempty"`
	RootTaskID string   // identifier of the root task
}

// Task represents an individual task (binary operation).
type Task struct {
	ID            string   `json:"id"`
	ExpressionID  string   // the expression this task belongs to
	Operator      string   `json:"operation"`
	Arg1          *float64 `json:"arg1,omitempty"`
	Arg2          *float64 `json:"arg2,omitempty"`
	DepTask1      string   // if the first operand comes from another task
	DepTask2      string   // if the second operand comes from another task
	OperationTime int      `json:"operation_time"` // operation execution time in milliseconds
	Status        string   // "pending", "running", "done"
	Result        *float64 `json:"result,omitempty"`
}

// Node represents a node in the expression tree.
type Node struct {
	IsLiteral bool
	Value     float64 // if IsLiteral is true
	Operator  string  // if node represents an operation
	Left      *Node
	Right     *Node
	TaskID    string // identifier of the corresponding task
}

// getOperationTime returns the operation execution time using environment variables.
func getOperationTime(op string) int {
	var envVar string
	switch op {
	case "+":
		envVar = "TIME_ADDITION_MS"
	case "-":
		envVar = "TIME_SUBTRACTION_MS"
	case "*":
		envVar = "TIME_MULTIPLICATIONS_MS"
	case "/":
		envVar = "TIME_DIVISIONS_MS"
	default:
		return 1000
	}
	msStr := os.Getenv(envVar)
	if msStr == "" {
		return 1000
	}
	ms, err := strconv.Atoi(msStr)
	if err != nil {
		return 1000
	}
	return ms
}

// buildExpressionTree builds an expression tree from tokens in Reverse Polish Notation.
func buildExpressionTree(tokens []calculator.Token) (*Node, error) {
	var stack []*Node
	for _, token := range tokens {
		if token.IsOperand {
			val, _ := token.GetOperand()
			stack = append(stack, &Node{IsLiteral: true, Value: val})
		} else if token.IsOperator {
			if len(stack) < 2 {
				return nil, errors.New("not enough operands")
			}
			// Pop the right then the left operand
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			node := &Node{
				IsLiteral: false,
				Operator:  token.Value.(string),
				Left:      left,
				Right:     right,
			}
			stack = append(stack, node)
		} else {
			return nil, errors.New("unexpected token")
		}
	}
	if len(stack) != 1 {
		return nil, errors.New("invalid expression")
	}
	return stack[0], nil
}

// createTasksFromNode recursively creates tasks from the expression tree.
// If the node represents an operation, a task is generated and its identifier is returned.
func createTasksFromNode(exprID string, node *Node) string {
	if node.IsLiteral {
		return ""
	}
	var dep1, dep2 string
	if !node.Left.IsLiteral {
		dep1 = createTasksFromNode(exprID, node.Left)
	}
	if !node.Right.IsLiteral {
		dep2 = createTasksFromNode(exprID, node.Right)
	}
	var arg1, arg2 *float64
	if node.Left.IsLiteral {
		arg1 = &node.Left.Value
	}
	if node.Right.IsLiteral {
		arg2 = &node.Right.Value
	}
	task := &Task{
		ID:            uuid.New().String(),
		ExpressionID:  exprID,
		Operator:      node.Operator,
		Arg1:          arg1,
		Arg2:          arg2,
		DepTask1:      dep1,
		DepTask2:      dep2,
		OperationTime: getOperationTime(node.Operator),
		Status:        "pending",
	}
	node.TaskID = task.ID
	storeMutex.Lock()
	tasksStore[task.ID] = task
	storeMutex.Unlock()
	return task.ID
}

// BuildExpressionTasks accepts an expression string, builds the tree, and generates tasks.
func BuildExpressionTasks(expression string) (*Expression, error) {
	tokens, err := calculator.Tokenize(expression)
	if err != nil {
		return nil, err
	}
	rpn, err := calculator.ShuntingYard(tokens)
	if err != nil {
		return nil, err
	}
	tree, err := buildExpressionTree(rpn)
	if err != nil {
		return nil, err
	}
	exprID := uuid.New().String()
	expr := &Expression{
		ID:     exprID,
		Expr:   expression,
		Status: "pending",
	}
	if tree.IsLiteral {
		expr.Status = "done"
		expr.Result = &tree.Value
	} else {
		rootTaskID := createTasksFromNode(exprID, tree)
		expr.RootTaskID = rootTaskID
	}
	storeMutex.Lock()
	expressionsStore[exprID] = expr
	storeMutex.Unlock()
	return expr, nil
}

// updateTaskDependencies checks whether the task's dependencies are ready and assigns their results.
func updateTaskDependencies(task *Task) bool {
	if task.DepTask1 != "" && task.Arg1 == nil {
		depTask, exists := tasksStore[task.DepTask1]
		if !exists || depTask.Status != "done" {
			return false
		}
		task.Arg1 = depTask.Result
	}
	if task.DepTask2 != "" && task.Arg2 == nil {
		depTask, exists := tasksStore[task.DepTask2]
		if !exists || depTask.Status != "done" {
			return false
		}
		task.Arg2 = depTask.Result
	}
	return task.Arg1 != nil && task.Arg2 != nil
}
