import { KLOROFORM_NAMESPACES, KLOROFORM_IGNORED_NAMESPACES } from "./config"
import { coreV1Api, appsV1Api } from "./k8s"

const getAllNamespaces = async () => {
  try {
    const res = await coreV1Api.listNamespace()
    return res.body.items.map((ns) => ns.metadata?.name) as string[]
  } catch (err) {
    console.error(err)
  }
}

export const getNamespaces = async () => {
  if (KLOROFORM_NAMESPACES)
    return KLOROFORM_NAMESPACES?.split(",").map((ns) => ns.trim())
  const namespaces = await getAllNamespaces()
  if (!namespaces) throw "No namespaces"

  return namespaces.filter(
    (ns) =>
      !KLOROFORM_IGNORED_NAMESPACES.split(",")
        .map((ns) => ns.trim())
        .includes(ns)
  )
}

export const getDeploymentsOfNamespace = async (namespace: string) => {
  try {
    const res = await appsV1Api.listNamespacedDeployment(namespace)
    return res.body.items
  } catch (err) {
    console.error(err)
  }
}
