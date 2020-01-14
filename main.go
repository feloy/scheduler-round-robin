package main

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	schedulerName = "scheduler-round-robin"
)

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Get nodes
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var nodeNames []string
	for _, node := range nodes.Items {
		isMaster := false
		for k := range node.Labels {
			if k == "node-role.kubernetes.io/master" {
				isMaster = true
				break
			}
		}
		if !isMaster {
			nodeNames = append(nodeNames, node.Name)
		}
	}
	fmt.Printf("found %d nodes: %v\n", len(nodeNames), nodeNames)

	watch, err := clientset.CoreV1().Pods("").Watch(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.schedulerName=%s,spec.nodeName=", schedulerName),
	})
	if err != nil {
		panic(err.Error())
	}

	nodeToSchedule := 0
	for event := range watch.ResultChan() {
		if event.Type != "ADDED" {
			continue
		}
		p := event.Object.(*v1.Pod)
		fmt.Println("found a pod to schedule:", p.Namespace, "/", p.Name)

		nodeName := nodeNames[nodeToSchedule]
		nodeToSchedule = (nodeToSchedule + 1) % len(nodeNames)

		err = clientset.CoreV1().Pods(p.Namespace).Bind(&v1.Binding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      p.Name,
				Namespace: p.Namespace,
			},
			Target: v1.ObjectReference{
				APIVersion: "v1",
				Kind:       "Node",
				Name:       nodeName,
			},
		})
		if err != nil {
			panic(err.Error())
		}

		timestamp := time.Now().UTC()
		_, err = clientset.CoreV1().Events(p.Namespace).Create(&v1.Event{
			Count:          1,
			Message:        fmt.Sprintf("pod %s scheduled to node %s", p.Name, nodeName),
			Reason:         "Scheduled",
			LastTimestamp:  metav1.NewTime(timestamp),
			FirstTimestamp: metav1.NewTime(timestamp),
			Type:           "Normal",
			Source: v1.EventSource{
				Component: schedulerName,
			},
			InvolvedObject: v1.ObjectReference{
				Kind:      "Pod",
				Name:      p.Name,
				Namespace: p.Namespace,
				UID:       p.UID,
			},
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: p.Name + "-",
			},
		})
		if err != nil {
			panic(err.Error())
		}
	}
}
