import dotenv from "dotenv"
dotenv.config()

export const {
  KLOROFORM_ANNOTATION_KEY = "kloroform/original-replica-count",
  KLOROFORM_NAMESPACES,
  KLOROFORM_IGNORED_NAMESPACES = "kube-system,kube-public,kube-node-lease,longhorn-system,cnpg-system,monitoring",
} = process.env
