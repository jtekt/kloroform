import * as k8s from "@kubernetes/client-node"

const kc = new k8s.KubeConfig()
kc.loadFromDefault()

const coreV1Api = kc.makeApiClient(k8s.CoreV1Api)
const appsV1Api = kc.makeApiClient(k8s.AppsV1Api)

const getNamespaces = async () => {
  try {
    const res = await coreV1Api.listNamespace()
    return res.body.items
  } catch (err) {
    console.error(err)
  }
}

const getDeploymentsOfNamespace = async (namespace: string) => {
  try {
    const res = await appsV1Api.listNamespacedDeployment(namespace)
    return res.body.items
  } catch (err) {
    console.error(err)
  }
}

const namespace = "oohara"

const main = async () => {
  // const namespaces = await getNamespaces()
  const deployments = await getDeploymentsOfNamespace(namespace)
  // const deloymentNames = deployments?.map((d) => d.metadata?.name)
  if (!deployments) throw "Failed to get deployments"
  for await (const deployment of deployments) {
    try {
      const deploymentName = deployment.metadata?.name
      if (!deploymentName) throw "Deployment has no name"
      if (!deployment.spec) throw "Deployment has no replicas"
      deployment.spec.replicas = 1

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

main()
