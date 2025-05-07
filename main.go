package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)




func main() {

	fmt.Print("= Kloroform =\n")

	kloroformAnnotationKey := "kloroform/original-replica-count"
	baseExceptions := "kube-system,kube-public,kube-node-lease,longhorn-system,cnpg-system"
	home := homedir.HomeDir()

	// Loading Kubeconfig
	var kubeconfig *string
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	
	wakeFlag := flag.Bool("wake", false, "Wake the cluster up")
	namespacesFlag := flag.String("namespaces", "", "(optional) comma separated list of namespaces to sedate/wake")
	exceptionsFlag := flag.String("exceptions", "", "(optional) comma separated list of namespaces to be excluded")
	
	
	flag.Parse()
	
	if(*wakeFlag) {
		fmt.Print("Operation: wake\n")
	} else {
		fmt.Print("Operation: sedate\n")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ignoredNamespace := baseExceptions + "," + *exceptionsFlag
	namespaces := listNamespaces(namespacesFlag, clientset, ignoredNamespace)
	
	for i, ns := range namespaces {

		fmt.Printf("- Namespace %d of %d: %s\n", i, len(namespaces), ns)

		deployments := listDeployments(clientset, ns)

		for _,deployment := range(deployments.Items) {
			fmt.Printf("  - Deployment: %s: ", deployment.Name)

			if *wakeFlag {
				wakeDeployment(clientset, ns, deployment, kloroformAnnotationKey)
			} else {
				sedateDeployment(clientset, ns, deployment, kloroformAnnotationKey)
			}
		}


		statefulSets := listStatefulSets(clientset, ns)

		for _,statefulSet := range(statefulSets.Items) {
			fmt.Printf("  - StatefulSet: %s: ", statefulSet.Name)

			if *wakeFlag {
				wakeStatefulSet(clientset, ns, statefulSet, kloroformAnnotationKey)
			} else {
				sedateStatefulSet(clientset, ns, statefulSet, kloroformAnnotationKey)
			}
		}
	}
}