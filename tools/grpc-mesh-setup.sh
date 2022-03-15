
function meshSetupIstio() {
  gcloud container hub memberships list --uri
  gcloud container hub memberships register MEMBERSHIP_NAME \
       --gke-cluster=GKE_CLUSTER \
       --enable-workload-identity \
       --project PROJECT_ID

  gcloud container hub mesh enable
  gcloud container clusters update CLUSTER_NAME --zone ZONE\
       --update-labels mesh_id=proj-PROJECT_NUMBER

  gcloud alpha container hub mesh update \
     --control-plane automatic \
     --membership $CLUSTER_NAME

  gcloud alpha container hub mesh describe

  kpt pkg get https://github.com/GoogleCloudPlatform/anthos-service-mesh-packages.git/asm@release-1.12 asm
  kubectl delete -f asm/canonical-service/controller.yaml

}
