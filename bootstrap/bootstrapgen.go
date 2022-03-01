package bootstrap

import (
	"errors"
	"io/ioutil"
	"os"
)

// channel_creds - used to authenticate XDS conn, google_default only, using transport credentials
// certs: mesh_ca (credentials/tls/cert_provider/meshca, file_watcher

const template = `
{
  "xds_servers": [
    {
      "server_uri": "localhost:15010",
      "channel_creds": [{"type": "insecure"}],
      "server_features" : ["xds_v3"]
    }
  ],
  "server_listener_resource_name_template": "istiod.istio-system.svc.cluster.local:15012",
  "certificate_providers" : {
			 "files": {
				 "plugin_name": "file_watcher",
				 "config": {			
						"certificate_file":   "/a/b/cert.pem",
						"private_key_file":    "/a/b/key.pem",
						"ca_certificate_file": "/a/b/ca.pem",
						"refresh_interval":   "200s"
          }
				},
			 "mesh_ca": {
				 "plugin_name": "mesh_ca",
				 "config": {		
						"server": {
									"api_type": 2,
									"grpc_services": [
										{
											"googleGrpc": {
												"call_credentials": [
													{
														"sts_service": {
															"subject_token_path": "test-subject-token-path"
														}
													}
												]
											},
											"timeout": "10s"
										}
									]
								}
				}
			 }
		 },
    "node": {
    "id": "sidecar~10.0.0.1~bob.grpcprobe~ns.cluster.local",
    "metadata": {
      "GENERATOR": "grpc",
      "NAMESPACE": "grpcprobe",
      "LABELS": {"version":  "v1"},
      "INSTANCE_IPS": ["10.0.0.1"],
      "CLUSTER_ID": "grpcprobe"
    }
  }
}
`

// GenerateBootstrap will write a Istio bootstrap file in the location expected by gRPC, using
// Istio environment variables:
//
// XDS_ADDR - the address of the XDS server, defaults to istiod.istio-system.svc:15010 if cert not set, and 15012 if root cert found
// POD_NAMESPACE, LABELS - based on standard mounts
// ISTIO_META_env variables used like in regular Istio
// ...
func GenerateBootstrap() error {
	bootF := os.Getenv("GRPC_XDS_BOOTSTRAP")
	if bootF == "" {
		return errors.New("missing GRPC_XDS_BOOTSTRAP")
	}

	if _, err := os.Stat(bootF); true || os.IsNotExist(err) {
		// TODO: write the bootstrap file.
		err := ioutil.WriteFile(bootF, []byte(template), 0700)
		if err != nil {
			return err
		}
	}


	return nil

}
