import { KLOROFORM_ANNOTATION_KEY } from "./config"
import { appsV1Api } from "./k8s"
import { getDeploymentsOfNamespace, getNamespaces } from "./utils"

export const sedate = async () => {
  const namespaces = await getNamespaces()

  for await (const namespace of namespaces) {
    console.log(`Sedating deployments of namespace ${namespace}`)
    const deployments = await getDeploymentsOfNamespace(namespace)
    if (!deployments) throw "Failed to get deployments"
    for await (const deployment of deployments) {
      try {
        if (!deployment.metadata?.name) throw "Deployment has no name"
        if (!deployment.spec) throw "Deployment has no spec"

        const deploymentName = deployment.metadata.name

        console.log(`Sedating deployment ${deploymentName}`)

        if (deployment.spec.replicas === 0)
          throw `Deployment ${deploymentName} already has 0 replicas`

        // Update replicas
        deployment.metadata.annotations = {
          ...deployment.metadata.annotations,
          [KLOROFORM_ANNOTATION_KEY]: String(deployment.spec.replicas),
        }
        deployment.spec.replicas = 0

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

sedate()
