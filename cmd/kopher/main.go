package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kitagry/gopherrotate"
	corev1 "k8s.io/api/core/v1"
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

	args := flag.Args()
	if len(args) < 1 {
		log.Println("arguments should be set")
		os.Exit(1)
	}

	mascot := gopherrotate.NewMascot()
	go func() {
		for {
			func() {
				podLog, err := clientset.CoreV1().Pods("default").GetLogs(args[0], &corev1.PodLogOptions{}).Stream(context.Background())
				if err != nil {
					fmt.Println(err)
					return
				}
				defer podLog.Close()

				buf := new(bytes.Buffer)
				_, err = io.Copy(buf, podLog)
				if err != nil {
					fmt.Println(err)
					return
				}
				mascot.Say(buf.String())
			}()
			time.Sleep(1 * time.Second)
		}
	}()

	err = mascot.Run()
	if err != nil {
		panic(err.Error())
	}
}
