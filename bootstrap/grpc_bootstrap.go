// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bootstrap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

	"google.golang.org/protobuf/types/known/structpb"

	log "google.golang.org/grpc/grpclog"
)

// Extracted from istio/istio/pkg/istio-agent/grpcxds
// - removed deps to istio internal packages

const (
	ServerListenerNamePrefix = "xds.istio.io/grpc/lds/inbound/"
	// ServerListenerNameTemplate for the name of the Listener resource to subscribe to for a gRPC
	// server. If the token `%s` is present in the string, all instances of the
	// token will be replaced with the server's listening "IP:port" (e.g.,
	// "0.0.0.0:8080", "[::]:8080").
	ServerListenerNameTemplate = ServerListenerNamePrefix + "%s"
)

// Bootstrap contains the general structure of what's expected by GRPC's XDS implementation.
// See https://github.com/grpc/grpc-go/blob/master/xds/internal/xdsclient/bootstrap/bootstrap.go
// TODO use structs from gRPC lib if created/exported
type Bootstrap struct {
	XDSServers                 []XdsServer                    `json:"xds_servers,omitempty"`
	Node                       *core.Node                   `json:"node,omitempty"`
	CertProviders              map[string]CertificateProvider `json:"certificate_providers,omitempty"`
	ServerListenerNameTemplate string                         `json:"server_listener_resource_name_template,omitempty"`
}



type ChannelCreds struct {
	Type   string      `json:"type,omitempty"`
	Config interface{} `json:"config,omitempty"`
}

type XdsServer struct {
	ServerURI      string         `json:"server_uri,omitempty"`
	ChannelCreds   []ChannelCreds `json:"channel_creds,omitempty"`
	ServerFeatures []string       `json:"server_features,omitempty"`
}

type CertificateProvider struct {
	PluginName string      `json:"plugin_name,omitempty"`
	Config     interface{} `json:"config,omitempty"`
}

func (cp *CertificateProvider) UnmarshalJSON(data []byte) error {
	var dat map[string]*json.RawMessage
	if err := json.Unmarshal(data, &dat); err != nil {
		return err
	}
	*cp = CertificateProvider{}

	if pluginNameVal, ok := dat["plugin_name"]; ok {
		if err := json.Unmarshal(*pluginNameVal, &cp.PluginName); err != nil {
			log.Warningf("failed parsing plugin_name in certificate_provider: %v", err)
		}
	} else {
		log.Warningf("did not find plugin_name in certificate_provider")
	}

	if configVal, ok := dat["config"]; ok {
		var err error
		switch cp.PluginName {
		case FileWatcherCertProviderName:
			config := FileWatcherCertProviderConfig{}
			err = json.Unmarshal(*configVal, &config)
			cp.Config = config
		default:
			config := FileWatcherCertProviderConfig{}
			err = json.Unmarshal(*configVal, &config)
			cp.Config = config
		}
		if err != nil {
			log.Warningf("failed parsing config in certificate_provider: %v", err)
		}
	} else {
		log.Warningf("did not find config in certificate_provider")
	}

	return nil
}

const FileWatcherCertProviderName = "file_watcher"

type FileWatcherCertProviderConfig struct {
	CertificateFile   string          `json:"certificate_file,omitempty"`
	PrivateKeyFile    string          `json:"private_key_file,omitempty"`
	CACertificateFile string          `json:"ca_certificate_file,omitempty"`
	RefreshDuration   string          `json:"refresh_interval,omitempty"`
}

func (c *FileWatcherCertProviderConfig) FilePaths() []string {
	return []string{c.CertificateFile, c.PrivateKeyFile, c.CACertificateFile}
}

// FileWatcherProvider returns the FileWatcherCertProviderConfig if one exists in CertProviders
func (b *Bootstrap) FileWatcherProvider() *FileWatcherCertProviderConfig {
	if b == nil || b.CertProviders == nil {
		return nil
	}
	for _, provider := range b.CertProviders {
		if provider.PluginName == FileWatcherCertProviderName {
			cfg, ok := provider.Config.(FileWatcherCertProviderConfig)
			if !ok {
				return nil
			}
			return &cfg
		}
	}
	return nil
}

// LoadBootstrap loads a Bootstrap from the given file path.
func LoadBootstrap(file string) (*Bootstrap, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	b := &Bootstrap{}
	if err := json.Unmarshal(data, b); err != nil {
		return nil, err
	}
	return b, err
}

type GenerateBootstrapOptions struct {
	// Original: Node             *model.Node
	ID string

	NodeMetadata     map[string]interface{}
	XdsUdsPath       string
	DiscoveryAddress string
	CertDir          string
	Locality         *core.Locality
}

// GenerateBootstrap generates the bootstrap structure for gRPC XDS integration.
func GenerateBootstrap(opts GenerateBootstrapOptions) (*Bootstrap, error) {
	xdsMeta, err := extractMeta(opts.NodeMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed extracting xds metadata: %v", err)
	}

	// TODO direct to CP should use secure channel (most likely JWT + TLS, but possibly allow mTLS)
	serverURI := opts.DiscoveryAddress
	if opts.XdsUdsPath != "" {
		serverURI = fmt.Sprintf("unix:///%s", opts.XdsUdsPath)
	}

	bootstrap := Bootstrap{
		XDSServers: []XdsServer{{
			ServerURI: serverURI,
			// connect locally via agent
			ChannelCreds:   []ChannelCreds{{Type: "insecure"}},
			ServerFeatures: []string{"xds_v3"},
		}},
		Node: &core.Node {
			Id:       opts.ID,
			Locality: opts.Locality,
			Metadata: xdsMeta,
		},
		ServerListenerNameTemplate: ServerListenerNameTemplate,
	}

	if opts.CertDir != "" {
		bootstrap.CertProviders = map[string]CertificateProvider{
			"default": {
				PluginName: "file_watcher",
				Config: FileWatcherCertProviderConfig{
					PrivateKeyFile:    path.Join(opts.CertDir, "key.pem"),
					CertificateFile:   path.Join(opts.CertDir, "cert-chain.pem"),
					CACertificateFile: path.Join(opts.CertDir, "root-cert.pem"),
					RefreshDuration:   "600s",
				},
			},
		}
	}

	return &bootstrap, err
}

func extractMeta(rawMeta map[string]interface{}) (*structpb.Struct, error) {
	xdsMeta, err := structpb.NewStruct(rawMeta)
	if err != nil {
		return nil, err
	}
	return xdsMeta, nil
}

// GenerateBootstrapFile generates and writes atomically as JSON to the given file path.
func GenerateBootstrapFile(opts GenerateBootstrapOptions, path string) (*Bootstrap, error) {
	bootstrap, err := GenerateBootstrap(opts)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.MarshalIndent(bootstrap, "", "  ")
	if err != nil {
		return nil, err
	}
	//if err := file.AtomicWrite(path, jsonData, os.FileMode(0o644)); err != nil {
	//	return nil, fmt.Errorf("failed writing to %s: %v", path, err)
	//}
	if err := ioutil.WriteFile(path, jsonData, os.FileMode(0o644)); err != nil {
		return nil, fmt.Errorf("failed writing to %s: %v", path, err)
	}
	return bootstrap, nil
}
