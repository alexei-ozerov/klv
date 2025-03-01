package mvc

import (
	"fmt"

	"github.com/alexei-ozerov/klv/kube"
	"github.com/alexei-ozerov/klv/tables"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"k8s.io/client-go/kubernetes"
)

// sessionState is used to track which model is focused
type sessionState uint

const (
	namespaceView  sessionState = 0
	podsView                    = 1
	containersView              = 2
	logsView                    = 3
)

type MainModel struct {
	state           sessionState
	index           int
	clientset       *kubernetes.Clientset
	namespace       string
	pod             string
	container       string
	namespaceTable  table.Model
	podsTable       table.Model
	containersTable table.Model
	logsTable       table.Model
	logSnippet      string
}

func NewModel() MainModel {
	m := MainModel{state: namespaceView}
	m.namespace = "default"
	m.pod = ""
	m.container = ""
	m.logSnippet = ""
	m.clientset = kube.InitKubeCtx()
	m.namespaceTable = tables.GetNamespacesTable(m.clientset)
	m.podsTable = tables.GetPodsTable(m.clientset, m.namespace)
	m.containersTable = tables.GetContainersTable(m.clientset, m.namespace, m.pod)
	m.logsTable = tables.GetLogsTable(m.clientset, m.namespace, m.pod, m.container)

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
				m.state = podsView
			} else {
				return m, tea.Quit
			}
		case "tab":
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
			if m.state == namespaceView {
				m.namespaceTable = tables.GetNamespacesTable(m.clientset)
			}
			if m.state == podsView {
				m.podsTable = tables.GetPodsTable(m.clientset, m.namespace)
			}
			if m.state == containersView {
				m.containersTable = tables.GetContainersTable(m.clientset, m.namespace, m.pod)
			}
			if m.state == logsView {
				m.logsTable = tables.GetLogsTable(m.clientset, m.namespace, m.pod, m.container)
			}
		case "enter":
			if m.state == namespaceView {
				m.namespace = fmt.Sprintf("%s", m.namespaceTable.SelectedRow()[0])
				m.podsTable = tables.GetPodsTable(m.clientset, m.namespace)
			}
			if m.state == podsView {
				if len(m.podsTable.SelectedRow()) > 0 {
					m.pod = fmt.Sprintf("%s", m.podsTable.SelectedRow()[0])
				}
				m.containersTable = tables.GetContainersTable(m.clientset, m.namespace, m.pod)
				m.state = containersView
			}
			if m.state == containersView {
				m.container = fmt.Sprintf("%s", m.containersTable.SelectedRow()[0])
				m.logsTable = tables.GetLogsTable(m.clientset, m.namespace, m.pod, m.container)
			}
			if m.state == logsView {
				logline := fmt.Sprintf("%s", m.logsTable.SelectedRow()[0])
				loglineWrapped := wrapText(logline, TextLength)
				m.logSnippet = loglineWrapped
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
			s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.containersTable.View()), logTableModelStyle.Render(m.logsTable.View())), focusedTextModelStyle.Render(fmt.Sprintf("%s", m.logSnippet)))
		} else { // log view
			s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(m.containersTable.View()), fLogModelStyle.Render(m.logsTable.View())), focusedTextModelStyle.Render(fmt.Sprintf("%s", m.logSnippet)))
		}
	} else {
		if m.state == namespaceView {
			s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.namespaceTable.View()), modelStyle.Render(m.podsTable.View()))
		} else {
			s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(m.namespaceTable.View()), focusedModelStyle.Render(m.podsTable.View()))
		}
	}

	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • j: scroll down • k: scroll up • enter: select item • r: reload table • q: exit\n"))

	return s
}

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

func wrapText(text string, width int) string {
	return wordwrap.String(text, width)
}
