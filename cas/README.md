# Integration with CAs 


## Istiod

Requires XDS address, JWT and CitadelRootCA.

## Google Managed Istio CA

Should work out of box for GKE projects.

## Google CA Service

https://cloud.google.com/traffic-director/docs/security-proxyless-setup

Workload certs issued have few additional infos. 1 and 3 seem some form of UUID or hex.
The '2' is more interesting - is the pod ID. This only appears to exist if they are created
by the node agent - not when creating them using the API ( which is reasonable since there 
is no way to verify, the tokens are not issued to pods.) It may be possible with the 
mounted JWTs - but I have not tried.

1.3.6.1.4.1.11129.2.6.1.1:
0&.$47a679c6-6362-49fc-a03a-b1ebf9e2b7a8

1.3.6.1.4.1.11129.2.6.1.2:
0...istiod-86f479d66c-d9mcd

1.3.6.1.4.1.11129.2.6.1.3:
0&.$1a170649-81b4-4302-b031-d82714e9b894


{iso(1) identified-organization(3) dod(6) internet(1) private(4) enterprise(1) google(11129) 2(2) 4(4) 2(2)}
/iso/identified-organization/dod/internet/private/enterprise/google/2/4/2


```yaml

apiVersion: security.cloud.google.com/v1
kind: WorkloadCertificateConfig
metadata:
  name: default
spec:
  # Required. The CA service that issues your certificates.
  certificateAuthorityConfig:
    certificateAuthorityServiceConfig:
      endpointURI: ISSUING_CA_POOL_URI

  # Required. The key algorithm to use. Choice of RSA or ECDSA.
  #
  # To maximize compatibility with various TLS stacks, your workloads
  # should use keys of the same family as your root and subordinate CAs.
  #
  # To use RSA, specify configuration such as:
  #   keyAlgorithm:
  #     rsa:
  #       modulusSize: 4096
  #
  # Currently, the only supported ECDSA curves are "P256" and "P384", and the only
  # supported RSA modulus sizes are 2048, 3072 and 4096.
  keyAlgorithm:
    rsa:
      modulusSize: 4096

  # Optional. Validity duration of issued certificates, in seconds.
  #
  # Defaults to 86400 (1 day) if not specified.
  validityDurationSeconds: 86400

  # Optional. Try to start rotating the certificate once this
  # percentage of validityDurationSeconds is remaining.
  #
  # Defaults to 50 if not specified.
  rotationWindowPercentage: 50
---
    apiVersion: security.cloud.google.com/v1
    kind: WorkloadCertificateConfig
    metadata:
      annotations:
        kubectl.kubernetes.io/last-applied-configuration: |
          {"apiVersion":"security.cloud.google.com/v1","kind":"WorkloadCertificateConfig","metadata":{"annotations":{},"name":"default"},"spec":{"certificateAuthorityConfig":{"certificateAuthorityServiceConfig":{"endpointURI":"//privateca.googleapis.com/projects/costin-asm1/locations/us-central1/caPools/mesh"}},"keyAlgorithm":{"rsa":{"modulusSize":4096}},"rotationWindowPercentage":50,"validityDurationSeconds":86400}}
      creationTimestamp: "2022-05-16T15:49:18Z"
      generation: 1
      name: default
      resourceVersion: "506873511"
      uid: 5ce0ed09-ea1e-4f76-a81f-11ef51bc20af
    spec:
      certificateAuthorityConfig:
        certificateAuthorityServiceConfig:
          endpointURI: //privateca.googleapis.com/projects/costin-asm1/locations/us-central1/caPools/mesh
      keyAlgorithm:
        rsa:
          modulusSize: 4096
      rotationWindowPercentage: 50
      validityDurationSeconds: 86400
    status:
      conditions:
        - lastTransitionTime: "2022-05-16T15:49:18Z"
          message: WorkloadCertificateConfig is ready
          observedGeneration: 1
          reason: ConfigReady
          status: "True"
          type: Ready

```

Annotation: 
```yaml
annotations:
  security.cloud.google.com/use-workload-certificates: ""

```

```Makefile

# Setup CAS and create the root CA
cas/setup:
	gcloud privateca pools create --project "${CONFIG_PROJECT_ID}" mesh --tier devops --location ${REGION}

	# Creates projects/PROJECT_ID/locations/LOCATION/caPools/mesh/certificateAuthorities/mesh-selfsigned
	# May want to use O=MESH_ID, for multi-project.
	# Google managed
	gcloud privateca roots create --project "${CONFIG_PROJECT_ID}" mesh-selfsigned --pool mesh --location ${REGION} \
		--auto-enable \
        --subject "CN=${PROJECT_ID}, O=${CONFIG_PROJECT_ID}"

	# In multi-project mode, workloads will still get a K8S token from the config project - which is exchanged with a certificate
	gcloud privateca pools --project "${CONFIG_PROJECT_ID}" add-iam-policy-binding mesh \
        --project "${CONFIG_PROJECT_ID}" \
        --location "${REGION}" \
        --member "group:${CONFIG_PROJECT_ID}.svc.id.goog:/allAuthenticatedUsers/" \
        --role "roles/privateca.workloadCertificateRequester"

# Setup the config cluster to use workload certificates.
cas/setup-cluster: CONFIG_PROJNUM=$(shell gcloud projects describe ${CONFIG_PROJECT_ID} --format="value(projectNumber)")
cas/setup-cluster:
	gcloud container clusters update ${CLUSTER_NAME} --project "${CONFIG_PROJECT_ID}" --region ${CLUSTER_LOCATION} --enable-mesh-certificates
	gcloud privateca pools add-iam-policy-binding --project "${CONFIG_PROJECT_ID}" mesh \
	  --location ${REGION} \
	  --role roles/privateca.certificateManager \
	  --member="serviceAccount:service-${CONFIG_PROJNUM}@container-engine-robot.iam.gserviceaccount.com"

cas/setup-k8s:
	cat manifests/cas-template.yaml | envsubst | kubectl apply -f -

#gcloud privateca pools add-iam-policy-binding istio \
#   --role=roles/privateca.workloadCertificateRequester
# --member="serviceAccount:service-601426346923@gcp-sa-meshdataplane.iam.gserviceaccount.com"
#--project wlhe-cr --location=us-central1

```

The pool is regional, for mesh we'll use the well-known name 'mesh' (customization possible but not a priority)

Any GSA requesting certificates must have the IAM binding. (roles/privateca.auditor may also be needed).

Certificate format:

```shell
/var/run/secrets/workload-spiffe-credentials$ ls
ca_certificates.pem  certificates.pem  private_key.pem
    
istio-proxy@istiod-86f479d66c-d9mcd:/var/run/secrets/workload-spiffe-credentials$ openssl x509 -text -in certificates.pem 
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            c1:c6:92:5e:0d:e2:d0:60:57:1a:be:6f:35:fb:94:18:74:16:1b
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: O = costin-asm1, CN = costin-asm1
        Validity
            Not Before: Jan 13 16:55:47 2023 GMT
            Not After : Jan 14 16:55:46 2023 GMT
        Subject: CN = 68dbfd00-0c6a-4dd0-ae93-454d3650b848
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (4096 bit)
                Modulus:
                    00:d6:7e:d0:03:76:34:04:0f:2b:9e:69:5b:f1:af:
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment, Key Agreement
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication, TLS Web Client Authentication
            X509v3 Basic Constraints: critical
                CA:FALSE
            X509v3 Subject Key Identifier: 
                3A:16:86:08:17:68:98:B1:98:B5:18:18:5F:AC:C8:76:93:96:EF:A2
            X509v3 Authority Key Identifier: 
                DB:91:C8:D5:9C:C6:9D:B3:73:F3:46:93:50:5E:3D:B5:63:98:9A:EB
            Authority Information Access: 
                CA Issuers - URI:http://privateca-content-62698b9b-0000-24ec-8817-001a113ab5ae.storage.googleapis.com/ce69187e15c38c33903f/ca.crt
            X509v3 Subject Alternative Name: 
                URI:spiffe://costin-asm1.svc.id.goog/ns/istio-system/sa/istiod
            1.3.6.1.4.1.11129.2.6.1.1: 
                0&.$47a679c6-6362-49fc-a03a-b1ebf9e2b7a8
            1.3.6.1.4.1.11129.2.6.1.2: 
                0...istiod-86f479d66c-d9mcd
            1.3.6.1.4.1.11129.2.6.1.3: 
                0&.$1a170649-81b4-4302-b031-d82714e9b894
    Signature Algorithm: sha256WithRSAEncryption
    Signature Value:
        0d:f1:cd:04:07:49:9b:37:02:ac:e1:5e:b0:dc:2e:1e:26:49:
-----BEGIN CERTIFICATE-----
MIIG4TCCBMmgAwIBAgIUAMHGkl4N4tBgVxq+bzX7lBh0FhswDQYJKoZIhvcNAQEL
-----END CERTIFICATE-----


istio-proxy@istiod-86f479d66c-d9mcd:/var/run/secrets/workload-spiffe-credentials$ openssl x509 -text -in ca_certificates.pem 
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            d6:71:a0:51:b6:70:d5:e3:ff:36:07:a6:de:7c:ab:49:88:03:6d
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: O = costin-asm1, CN = costin-asm1
        Validity
            Not Before: Apr 23 17:09:47 2022 GMT
            Not After : Apr 23 03:17:27 2032 GMT
        Subject: O = costin-asm1, CN = costin-asm1
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (4096 bit)
                Modulus:
                    00:b4:49:ad:8b:b6:5f:e5:97:21:2f:d1:63:3f:75:
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Certificate Sign, CRL Sign
            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Subject Key Identifier: 
                DB:91:C8:D5:9C:C6:9D:B3:73:F3:46:93:50:5E:3D:B5:63:98:9A:EB
            X509v3 Authority Key Identifier: 
                DB:91:C8:D5:9C:C6:9D:B3:73:F3:46:93:50:5E:3D:B5:63:98:9A:EB
    Signature Algorithm: sha256WithRSAEncryption
    Signature Value:
        28:38:69:9a:aa:3d:3c:1f:d1:cd:af:a4:47:1f:4a:03:bf:3b:
```
