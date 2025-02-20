package main

import "github.com/ZolotarevAlexandr/yl_sprint_2_final/agent/agent"

func main() {
	go agent.RunAgent()
	select {}
}
