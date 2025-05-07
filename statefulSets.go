package main

import (
	"context"
	"fmt"
	"strconv"

	appsV1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


func listStatefulSets(clientset *kubernetes.Clientset, ns string) *appsV1.StatefulSetList{
	statefulSets, err := clientset.AppsV1().StatefulSets(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
	return statefulSets
}


func sedateStatefulSet(clientset *kubernetes.Clientset, ns string,statefulSet appsV1.StatefulSet, kloroformAnnotationKey string){
	// Do nothing if already sedated
	if(*statefulSet.Spec.Replicas == 0) {
		fmt.Printf("Already has 0 replicas \n")
		return
	} 

	annotationsMap := statefulSet.GetAnnotations()
	annotationsMap[kloroformAnnotationKey] = strconv.Itoa(int(*statefulSet.Spec.Replicas))
	statefulSet.SetAnnotations(annotationsMap)
	*statefulSet.Spec.Replicas = 0

	clientset.AppsV1().StatefulSets(ns).Update(context.TODO(), &statefulSet, metav1.UpdateOptions{})

	fmt.Printf("Scaled down to 0 replicas\n")
}


func wakeStatefulSet(clientset *kubernetes.Clientset, ns string,statefulSet appsV1.StatefulSet, kloroformAnnotationKey string){
	annotationsMap := statefulSet.GetAnnotations()

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

	
	*statefulSet.Spec.Replicas = int32(originalReplicaCount)

	clientset.AppsV1().StatefulSets(ns).Update(context.TODO(), &statefulSet, metav1.UpdateOptions{})
	fmt.Printf("Scaled up to %d replicas\n", originalReplicaCount)
}
