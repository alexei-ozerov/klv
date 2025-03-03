package mvc

import (
	"fmt"

	"github.com/alexei-ozerov/klv/kube"
	"github.com/alexei-ozerov/klv/tables"
	"github.com/alexei-ozerov/klv/utils"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/kubernetes"
)

// sessionState is used to track which model is focused
type sessionState uint

const TextLength = 150

var (
	modelStyle = lipgloss.NewStyle().
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.HiddenBorder())

	focusedModelStyle = lipgloss.NewStyle().
				Align(lipgloss.Left, lipgloss.Left).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("69"))

	fLogModelStyle = lipgloss.NewStyle().
			Align(lipgloss.Left, lipgloss.Left).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("69"))

	logTableModelStyle = lipgloss.NewStyle().
				Align(lipgloss.Left, lipgloss.Left).
				BorderStyle(lipgloss.HiddenBorder())

	focusedTextModelStyle = lipgloss.NewStyle().
				Height(5).
				Width(TextLength).
				Align(lipgloss.Left, lipgloss.Left).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("69"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

const (
	namespaceView  sessionState = 0
	podsView                    = 1
	containersView              = 2
	logsView                    = 3
)

type MainModel struct {
	state             sessionState
	clientset         *kubernetes.Clientset
	selectedNamespace string
	selectedPod       string
	selectedContainer string
	selectedLogLine   string
	namespaceTable    table.Model
	podsTable         table.Model
	containersTable   table.Model
	logsTable         table.Model
}

type ErrorMsg struct {
	Error string
}

func NewModel() MainModel {
	m := MainModel{state: namespaceView}
	m.selectedNamespace = "default"
	m.selectedPod = ""
	m.selectedContainer = ""
	m.selectedLogLine = ""
	m.clientset = kube.InitKubeCtx()
	m.namespaceTable = tables.GetNamespacesTable(m.clientset)
	m.podsTable = tables.GetPodsTable(m.clientset, m.selectedNamespace)
	m.containersTable = tables.GetContainersTable(m.clientset, m.selectedNamespace, m.selectedPod)
	m.logsTable = tables.GetLogsTable(m.clientset, m.selectedNamespace, m.selectedPod, m.selectedContainer)

	return m
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == containersView || m.state == logsView {
				m.selectedLogLine = ""
				m.logsTable = tables.ClearLogsTable()
				m.state = podsView
			} else {
				return m, tea.Quit
			}
		case "tab", "h":
			if m.state == namespaceView || m.state == podsView {
				if m.state == namespaceView {
					m.state = podsView
				} else {
					m.state = namespaceView
				}
			} else {
				if m.state == containersView {
					m.state = logsView
				} else {
					m.state = containersView
				}
			}
		case "r":
			switch m.state {
			case namespaceView:
				m.namespaceTable = tables.GetNamespacesTable(m.clientset)
			case podsView:
				m.podsTable = tables.GetPodsTable(m.clientset, m.selectedNamespace)
			case containersView:
				m.containersTable = tables.GetContainersTable(m.clientset, m.selectedNamespace, m.selectedPod)
			case logsView:
				m.logsTable = tables.GetLogsTable(m.clientset, m.selectedNamespace, m.selectedPod, m.selectedContainer)
				m.logsTable.GotoBottom()
			}
		case "enter", "l":
			switch m.state {
			case namespaceView:
				m.selectedNamespace = fmt.Sprintf("%s", m.namespaceTable.SelectedRow()[0])
				m.podsTable = tables.GetPodsTable(m.clientset, m.selectedNamespace)
				m.state = podsView
			case podsView:
				// If more than one pod, do the thing :3
				if len(m.podsTable.SelectedRow()) > 0 {
					m.selectedPod = fmt.Sprintf("%s", m.podsTable.SelectedRow()[0])
					m.containersTable = tables.GetContainersTable(m.clientset, m.selectedNamespace, m.selectedPod)

					// Preload first container's log data
					m.selectedContainer = fmt.Sprintf("%s", m.containersTable.SelectedRow()[0])
					m.logsTable = tables.GetLogsTable(m.clientset, m.selectedNamespace, m.selectedPod, m.selectedContainer)
					m.logsTable.GotoBottom()

					// If selection is valid, change screen focus
					if m.selectedPod != "" {
						m.state = containersView
					}
				}
			case containersView:
				m.selectedContainer = fmt.Sprintf("%s", m.containersTable.SelectedRow()[0])
				m.logsTable = tables.GetLogsTable(m.clientset, m.selectedNamespace, m.selectedPod, m.selectedContainer)
				m.logsTable.GotoBottom()
				if m.selectedContainer != "" {
					m.state = logsView
				}
			case logsView:
				logline := fmt.Sprintf("%s", m.logsTable.SelectedRow()[0])
				loglineWrapped := utils.WrapText(logline, TextLength)
				m.selectedLogLine = loglineWrapped
			}
		}
		switch m.state {
		case logsView:
			m.logsTable, cmd = m.logsTable.Update(msg)
		case containersView:
			m.containersTable, cmd = m.containersTable.Update(msg)
		case podsView:
			m.podsTable, cmd = m.podsTable.Update(msg)
		case namespaceView:
			m.namespaceTable, cmd = m.namespaceTable.Update(msg)
		}
	}
	return m, tea.Batch(cmd)
}

func (m MainModel) View() string {
	var s string

	if m.state == logsView || m.state == containersView {
		if m.state == containersView {
			s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.containersTable.View()), logTableModelStyle.Render(m.logsTable.View())), focusedTextModelStyle.Render(fmt.Sprintf("%s", m.selectedLogLine)))
		} else { // log view
			s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(m.containersTable.View()), fLogModelStyle.Render(m.logsTable.View())), focusedTextModelStyle.Render(fmt.Sprintf("%s", m.selectedLogLine)))
		}
	} else {
		if m.state == namespaceView {
			s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.namespaceTable.View()), modelStyle.Render(m.podsTable.View()))
		} else {
			s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(m.namespaceTable.View()), focusedModelStyle.Render(m.podsTable.View()))
		}
	}

	s += helpStyle.Render(fmt.Sprintf("\ntab, h: cycle next • j: scroll down • k: scroll up • enter, l: select item • r: reload table • q: exit\n"))

	return s
}
