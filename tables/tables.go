package tables

import (
	"fmt"
//	"slices"

	"github.com/alexei-ozerov/klv/kube"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/kubernetes"
)

func GetNamespacesTable(clientset *kubernetes.Clientset) table.Model {
	columns := []table.Column{
		{Title: "Namespace", Width: 35},
	}

	ns := kube.GetNamespace(clientset)

	var rows []table.Row
	for _, namespace := range ns {
		rows = append(rows, table.Row{namespace})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(false).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func GetPodsTable(clientset *kubernetes.Clientset, ns string) table.Model {
	columns := []table.Column{
		{Title: fmt.Sprintf("Pods: %s", ns), Width: 70},
	}

	pods := kube.GetPods(clientset, ns)

	var rows []table.Row
	for _, pod := range pods {
		rows = append(rows, table.Row{pod})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(false).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func GetContainersTable(clientset *kubernetes.Clientset, ns, pod string) table.Model {
	columns := []table.Column{
		{Title: fmt.Sprintf("Containers: %s", ns), Width: 30},
	}

	var rows []table.Row
	if pod != "" {
		containers := kube.GetContainers(clientset, ns, pod)

		for _, container := range containers {
			rows = append(rows, table.Row{container})
		}
	} else {
		rows = append(rows, table.Row{""})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(false).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func GetLogsTable(clientset *kubernetes.Clientset, ns, pod, container string) table.Model {
	columns := []table.Column{
		{Title: fmt.Sprintf("Logs: %s - %s", pod, container), Width: 114},
	}

	var rows []table.Row
	if container != "" {
		logs := kube.GetLogs(clientset, ns, pod, container)
		//	slices.Reverse(logs)

		for _, log := range logs {
			rows = append(rows, table.Row{log})
		}
	} else {
		rows = append(rows, table.Row{""})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(false).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func ClearLogsTable() table.Model {
	columns := []table.Column{
		{Title: fmt.Sprintf("Logs: "), Width: 114},
	}

	var rows []table.Row
	rows = append(rows, table.Row{""})

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(false).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
