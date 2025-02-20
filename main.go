package main

import (
	"github.com/ZolotarevAlexandr/yl_sprint_2_final/agent/agent"
	"github.com/ZolotarevAlexandr/yl_sprint_2_final/orchestrator/orchestrator"
)

func main() {
	go orchestrator.RunOrchestrator()
	go agent.RunAgent()
	select {}
}
