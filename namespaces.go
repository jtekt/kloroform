package main

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)




func listNamespaces(namespacesFlag *string, clientset *kubernetes.Clientset, ignoredNamespace string) []string {
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