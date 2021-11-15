package main

import (
	"log"
	"math/rand"
	"strings"

	v1 "k8s.io/api/core/v1"
	extender "k8s.io/kube-scheduler/extender/v1"
)

const (
	// LuckyPred rejects a node if you're not lucky ¯\_(ツ)_/¯
	LuckyPred        = "Lucky"
	LuckyPredFailMsg = "Well, you're not lucky"
	BandPred = "Satisfy minimum bandwidth " + minBandwidth
	BandPredFailMsg = "Unsatisfy minimum bandwidth " + minBandwidth
)

// var predicatesFuncs = map[string]FitPredicate{
// 	LuckyPred: LuckyPredicate,
// }

var predicatesFuncs = map[string]FitPredicate{
	BandPred: BandPredicate,
}

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

// var predicatesSorted = []string{LuckyPred}
var predicatesSorted = []string{BandPred}

// filter filters nodes according to predicates defined in this extender
// it's webhooked to pkg/scheduler/core/generic_scheduler.go#findNodesThatFitPod()
func filter(args extender.ExtenderArgs) *extender.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(extender.FailedNodesMap)
	pod := args.Pod

	// TODO: parallelize this
	// TODO: handle error
	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}
	}

	result := extender.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	fits := true
	var failReasons []string
	for _, predicateKey := range predicatesSorted {
		fit, failures, err := predicatesFuncs[predicateKey](pod, node)
		if err != nil {
			return false, nil, err
		}
		fits = fits && fit
		failReasons = append(failReasons, failures...)
	}
	return fits, failReasons, nil
}

func LuckyPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	lucky := rand.Intn(2) == 0
	if lucky {
		log.Printf("pod %v/%v is lucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name)
		return true, nil, nil
	}
	log.Printf("pod %v/%v is unlucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name)
	return false, []string{LuckyPredFailMsg}, nil
}

func BandPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	// how to get node's bandwidth???
	// satisfy := (node.bandwidth >= minBandwidth)
	if satisfy {
		log.Printf("pod %v/%v is satisfied by minimum bandwidth " + minBandwidth " on node %v\n", pod.Name, pod.Namespace, node.Name)
		return true, nil, nil
	}
	log.Printf("pod %v/%v is unsatisfied by minimum bandwidth " + minBandwidth " on node %v\n", pod.Name, pod.Namespace, node.Name)
	return false, []string{BandPredFailMsg}, nil
}
