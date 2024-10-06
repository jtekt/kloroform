package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	appsV1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)


func listnamespaces(namespacesFlag *string, clientset *kubernetes.Clientset, ignoredNamespace string) []string {
	var namespaces []string

	if *namespacesFlag != "" {
		fmt.Printf("Namespaces (user provided): %s\n", *namespacesFlag)
		namespaces = strings.Split(*namespacesFlag, ",")
	} else {
		fmt.Printf("Namespaces: all except following: %s\n", ignoredNamespace)

		namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d namespaces in the cluster\n", len(namespaceList.Items))

		for _, ns := range namespaceList.Items {

			isException := false
			for _, exception := range strings.Split(ignoredNamespace, ",") {
				if exception == ns.Name {
					isException = true
				}
			}

			if(!isException){
				namespaces = append(namespaces, ns.Name)		
			}

		}
	}

	return namespaces
}

func sedateDeployment(clientset *kubernetes.Clientset, ns string,deployment appsV1.Deployment, kloroformAnnotationKey string){
	// Do nothing if already sedated
	if(*deployment.Spec.Replicas == 0) {
		fmt.Printf("Already has 0 replicas \n")
		return
	} 

								
	annotationsMap := deployment.GetAnnotations()
	annotationsMap[kloroformAnnotationKey] = strconv.Itoa(int(*deployment.Spec.Replicas))
	deployment.SetAnnotations(annotationsMap)
	*deployment.Spec.Replicas = 0

	clientset.AppsV1().Deployments(ns).Update(context.TODO(), &deployment, metav1.UpdateOptions{})

	fmt.Printf("Scaled down to 0 replicas\n")


	
}


func wakeDeployment(clientset *kubernetes.Clientset, ns string,deployment appsV1.Deployment, kloroformAnnotationKey string){
	annotationsMap := deployment.GetAnnotations()

	originalReplicaCountString := annotationsMap[kloroformAnnotationKey]

	if(originalReplicaCountString == "") {
		fmt.Printf("No kloroform annotation, skipping \n")
		return
	} 

	originalReplicaCount, err := strconv.ParseInt(originalReplicaCountString, 10, 32)

	if(err != nil) {
		panic(err)
	}

	delete(annotationsMap, kloroformAnnotationKey)

	
	*deployment.Spec.Replicas = int32(originalReplicaCount)

	clientset.AppsV1().Deployments(ns).Update(context.TODO(), &deployment, metav1.UpdateOptions{})
	fmt.Printf("Scaled up to %d replicas\n", originalReplicaCount)
	
}

func main() {

	fmt.Print("= Kloroform =\n")

	kloroformAnnotationKey := "kloroform/original-replica-count"
	baseExceptions := "kube-system,kube-public,kube-node-lease,longhorn-system,cnpg-system"

	// Loading Kubeconfig
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	
	wakeFlag := flag.Bool("wake", false, "Wake the cluster up")
	namespacesFlag := flag.String("namespaces", "", "comma separated list of namespaces to sedate/wake")
	exceptionsFlag := flag.String("exceptions", "", "comma separated list of namespaces to be excluded")
	
	
	flag.Parse()
	
	if(*wakeFlag) {
		fmt.Print("Operation: wake\n")
	} else {
		fmt.Print("Operation: sedate\n")
	}

	ignoredNamespace := baseExceptions + "," + *exceptionsFlag
	

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}


	namespaces := listnamespaces(namespacesFlag, clientset, ignoredNamespace)
	
	for i, ns := range namespaces {

		fmt.Printf("- Namespace %d of %d: %s\n", i, len(namespaces), ns)

		deployments, err := clientset.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		for _,deployment := range(deployments.Items) {
			fmt.Printf("  - Deployment: %s: ", deployment.Name)

			if *wakeFlag {
				wakeDeployment(clientset, ns, deployment, kloroformAnnotationKey)
			} else {
				sedateDeployment(clientset, ns, deployment, kloroformAnnotationKey)
			}
		}
	}
}