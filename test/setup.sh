
PROJECT_ID=${PROJECT_ID:-costin-asm1}
PROJECT_NUMBER=${PROJECT_NUMBER:-438684899409}

functio setupTD() {
  NAMESPACE=${1:-default}
  KSA=${2:-default}

  gcloud services enable trafficdirector.googleapis.com

  gcloud projects add-iam-policy-binding ${PROJECT_ID} \
      --member serviceAccount:${PROJECT_ID}.svc.id.goog[${NAMESPACE}/${KSA}] \
      --role=roles/trafficdirector.client

  gcloud projects add-iam-policy-binding ${PROJECT_ID} \
      --member serviceAccount:service-${PROJECT_NUMBER}@gcp-sa-meshdataplane.iam.gserviceaccount.com \
      --role=roles/trafficdirector.client
}

function setupWI() {
  NAMESPACE=${1:-default}
  KSA=${2:-default}

  GSA_NAME=k8s-${NAMESPACE}-${KSA}

  gcloud iam service-accounts create ${GSA_NAME} \
      --project=${PROJECT_ID}

#      gcloud projects add-iam-policy-binding ${PROJECT_ID} \
#          --member "serviceAccount:${GSA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
#          --role "ROLE_NAME"

  # KSA can impersonate GSA
  gcloud iam service-accounts add-iam-policy-binding k8s-fortio@costin-asm1.iam.gserviceaccount.com \
      --role roles/iam.workloadIdentityUser \
      --member "serviceAccount:${PROJECT_ID}.svc.id.goog[default/default]"


  gcloud iam service-accounts add-iam-policy-binding ${GSA_NAME}@${$PROJECT_ID}.iam.gserviceaccount.com \
      --role roles/iam.workloadIdentityUser \
      --member "serviceAccount:${PROJECT_ID}.svc.id.goog[${NAMESPACE}/${KSA}]"

 # Anotated KSA
 kubectl annotate serviceaccount ${KSA} \
     --namespace ${NAMESPACE} \
     iam.gke.io/gcp-service-account=${GSA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com
}

function getWIToken() {
  curl -H "Metadata-Flavor: Google" http://169.254.169.254/computeMetadata/v1/instance/service-accounts/

}
