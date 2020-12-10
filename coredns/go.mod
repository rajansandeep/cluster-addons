module sigs.k8s.io/cluster-addons/coredns

go 1.15

require (
	github.com/coredns/corefile-migration v1.0.10
	github.com/go-logr/logr v0.2.1
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.7.0-alpha.5
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20201209165851-b731a6217520
)
