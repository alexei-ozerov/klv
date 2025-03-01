package kube

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func InitKubeCtx() *kubernetes.Clientset {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "absolute path to the kubeconfig file")
	flag.Parse()

	// Build the config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset

}

func GetNamespace(clientset *kubernetes.Clientset) []string {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Print namespace names
	var nsList []string
	for _, namespace := range namespaces.Items {
		ns := fmt.Sprintf(namespace.Name)
		nsList = append(nsList, ns)
	}

	return nsList
}

func GetPods(clientset *kubernetes.Clientset, namespace string) []string {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Print namespace names
	var podList []string
	for _, pod := range pods.Items {
		podName := fmt.Sprintf(pod.Name)
		podList = append(podList, podName)
	}

	return podList
}

func GetContainers(clientset *kubernetes.Clientset, namespace string, podName string) []string {
	var containerNames []string

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		errorMessage := fmt.Sprintf("%s", err.Error())
		containerNames = append(containerNames, errorMessage)
	} else {
		for _, container := range pod.Spec.Containers {
			containerNames = append(containerNames, container.Name)
		}
	}

	return containerNames
}

func GetLogs(clientset *kubernetes.Clientset, namespace string, podName string, containerName string) []string {
	podLogOpts := &corev1.PodLogOptions{}
	if containerName != "" {
		podLogOpts.Container = containerName
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		panic(err.Error())
	}

	defer func(podLogs io.ReadCloser) {
		err := podLogs.Close()
		if err != nil {
			panic(err.Error())
		}
	}(podLogs)

	logs, err := io.ReadAll(podLogs)
	if err != nil {
		panic(err.Error())
	}

	logLines := strings.Split(string(logs), "\n")

	return logLines
}

// homeDir returns the home directory for the executing user.
func homeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return dirname
}
