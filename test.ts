import { getDeploymentsOfNamespace, getNamespaces } from "./utils"

export const test = async () => {
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

        console.log(`Would have sedated deployment ${deploymentName}`)
      } catch (error) {
        console.error(error)
      }
    }
  }
}

test()
