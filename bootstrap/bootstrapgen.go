package bootstrap

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"text/template"
)

// channel_creds - used to authenticate XDS conn, google_default only, using transport credentials
// certs: mesh_ca (credentials/tls/cert_provider/meshca, file_watcher

// TODO: if auto-mtls label is set, Istiod should generate MTLS config
// Label would be added by injector or user.
// This would match Istio behavior. Meanwhile - Policy.

// Server: "unix:///etc/istio/proxy/XDS" using agent, istiod.istio-system.svc:15010 for plaintext
const grpcTemplate = `
{
  "xds_servers": [
    {
      "server_uri": "${.Server}",
      "channel_creds": [{"type": "insecure"}],
      "server_features" : ["xds_v3"]
    }
  ],
  "node": {
    "id": "sidecar~${.IP}~${.Name}.${.Namespace}~${.Namespace}.cluster.local",
    "metadata": {
      "INSTANCE_IPS": "127.0.1.1",
      "PILOT_SAN": [
        "istiod.istio-system.svc"
      ],
      "GENERATOR": "grpc",
      "NAMESPACE": "${.Namespace}"
    },
    "localisty": {},
    "UserAgentVersionType": "istiov1"
  },
  "certificate_providers": {
    "default": {
      "plugin_name": "file_watcher",
      "config": {
        "certificate_file": "../../../../tests/testdata/certs/default/cert-chain.pem",
        "private_key_file": "../../../../tests/testdata/certs/default/key.pem",
        "ca_certificate_file": "../../../../tests/testdata/certs/default/root-cert.pem",
        "refresh_interval": "900s"
      }
    }
  },
  "server_listener_resource_name_template": "xds.istio.io/grpc/lds/inbound/%s"
}
`

//const oldTemplate = `
//{
//  "certificate_providers" : {
//			 "mesh_ca": {
//				 "plugin_name": "mesh_ca",
//				 "config": {
//						"server": {
//									"api_type": 2,
//									"grpc_services": [
//										{
//											"googleGrpc": {
//												"call_credentials": [
//													{
//														"sts_service": {
//															"subject_token_path": "test-subject-token-path"
//														}
//													}
//												]
//											},
//											"timeout": "10s"
//										}
//									]
//								}
//				}
//			 }
//		 },
//}
//`

func BootstrapData() []byte {
	// TODO: write the bootstrap file.
	t := template.New("grpc")
	_, err := t.Parse(grpcTemplate)
	if err != nil {
		return nil
	}
	out := &bytes.Buffer{}
	t.Execute(out, map[string]interface{}{})

	return out.Bytes()

}

// GenerateBootstrap will write a Istio bootstrap file in the location expected by gRPC, using
// Istio environment variables:
//
// XDS_ADDR - the address of the XDS server, defaults to istiod.istio-system.svc:15010 if cert not set, and 15012 if root cert found
// POD_NAMESPACE, LABELS - based on standard mounts
// ISTIO_META_env variables used like in regular Istio
// ...
func GenerateBootstrapTmpl() error {
	bootF := os.Getenv("GRPC_XDS_BOOTSTRAP")
	if bootF == "" {
		return errors.New("missing GRPC_XDS_BOOTSTRAP")
	}

	if _, err := os.Stat(bootF); os.IsNotExist(err) {
		// TODO: write the bootstrap file.
		err := ioutil.WriteFile(bootF, BootstrapData(), 0700)
		if err != nil {
			return err
		}
	}

	return nil

}
