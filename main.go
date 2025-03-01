package main

import (
	"log"

	"github.com/alexei-ozerov/klv/mvc"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(mvc.NewModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
