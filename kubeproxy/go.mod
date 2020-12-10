module addon-operators/kubeproxy

go 1.15

require (
	github.com/go-logr/logr v0.2.1
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.7.0-alpha.5
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20201209165851-b731a6217520
)
