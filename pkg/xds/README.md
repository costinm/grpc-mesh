Prototype minimal implementation of PTR-DS and related Istio protocols using connect-go.

This is intended for agent and small-size images, where the binary size is a concern but 
performance and advanced features are not.

For example, it can be used in an agent or minimal server running on a low-power device.

# Behavior

1. Look for certificates in the standard location. A CA agent must provision them
2. if no certificates - assume 'secure network'
3. use istiod.istio-system.svc by default, XDS_ADDR otherwise

