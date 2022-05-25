package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	projects := make(map[string]*exec.Cmd)
	for {
		podList, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("failed to get pods: %v\n", err)
			return
		}

		podNames := make([]string, len(podList.Items))
		for i, pod := range podList.Items {
			podNames[i] = pod.GetName()
		}

		for _, podName := range podNames {
			_, ok := projects[podName]
			if ok {
				continue
			}

			cmd := exec.CommandContext(ctx, "kopher", podName)
			cmd.Stdout = os.Stdout
			err := cmd.Start()
			if err != nil {
				fmt.Printf("failed to run command: %v\n", err)
			}
			projects[podName] = cmd
		}

		for name, cmd := range projects {
			if !in(podNames, name) {
				cmd.Process.Kill()
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func in(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
