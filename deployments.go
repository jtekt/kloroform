package main

import (
	"context"
	"fmt"
	"strconv"

	appsV1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


func listDeployments(clientset *kubernetes.Clientset, ns string) *appsV1.DeploymentList{
	deployments, err := clientset.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
	return deployments
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
