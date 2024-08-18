import { appsV1Api } from "./k8s"
import { getDeploymentsOfNamespace, getNamespaces } from "./utils"
import { KLOROFORM_ANNOTATION_KEY } from "./config"

export const wake = async () => {
  const namespaces = await getNamespaces()
  for await (const namespace of namespaces) {
    console.log(`Waking up deployments of namespace ${namespace}`)

    const deployments = await getDeploymentsOfNamespace(namespace)
    if (!deployments) throw "Failed to get deployments"
    for await (const deployment of deployments) {
      try {
        if (!deployment.metadata?.name) throw "Deployment has no name"
        if (!deployment.spec) throw "Deployment has no spec"

        const deploymentName = deployment.metadata.name
        console.log(`Waking up deployment ${deploymentName}`)

        if (deployment.spec.replicas !== 0)
          throw "Deployment already has replicas"

        let targetReplicaCount = 1

        if (
          deployment.metadata.annotations &&
          deployment.metadata.annotations[KLOROFORM_ANNOTATION_KEY]
        ) {
          const originalReplicaCount = Number(
            deployment.metadata.annotations[KLOROFORM_ANNOTATION_KEY]
          )

          if (!isNaN(originalReplicaCount)) {
            console.log(`Original replica count found: ${originalReplicaCount}`)
            targetReplicaCount = originalReplicaCount
          } else console.log(`Cannot parse original replica count`)

          delete deployment.metadata.annotations[KLOROFORM_ANNOTATION_KEY]
        } else {
          console.log(`Cannot find annotation for original replica count`)
        }

        deployment.spec.replicas = targetReplicaCount

        await appsV1Api.replaceNamespacedDeployment(
          deploymentName,
          namespace,
          deployment
        )
      } catch (error) {
        console.error(error)
      }
    }
  }
}

wake()
