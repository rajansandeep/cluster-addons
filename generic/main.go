/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	addonsv1alpha1 "sigs.k8s.io/cluster-addons/generic/api/v1alpha1"
	"sigs.k8s.io/cluster-addons/generic/controllers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = addonsv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var genericList addonsv1alpha1.GenericList

	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "generic-operator",
		Port:               9443,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	clientset, err := dynamic.NewForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	resourcesList, err := clientset.Resource(schema.GroupVersionResource{
		Group:    "addons.x-k8s.io",
		Version:  "v1alpha1",
		Resource: "generics",
	}).Namespace("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		setupLog.Error(err, "unable to get generic resource")
		os.Exit(1)
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(resourcesList.UnstructuredContent(), &genericList)
	if err != nil {
		setupLog.Error(err, "unable to destructure unstructured")
		os.Exit(1)
	}

	if len(genericList.Items) == 0 {
		fmt.Fprint(os.Stderr, "You need to define the addon that you want the controller to manager. Please create a"+
			" `Generic` resource for the addon\n")
		os.Exit(1)
	}

	// TODO: Add watch to Generic; respond to chanages
	for i := range genericList.Items {
		genericObject := genericList.Items[i].Spec

		gvk := schema.GroupVersionKind{
			Kind:    genericObject.ObjectKind.Kind,
			Version: genericObject.ObjectKind.Version,
			Group:   genericObject.ObjectKind.Group,
		}

		if err = (&controllers.GenericReconciler{
			Client:  mgr.GetClient(),
			GVK:     gvk,
			Log:     ctrl.Log.WithName("controllers").WithName(gvk.Kind),
			Scheme:  mgr.GetScheme(),
			Channel: genericObject.Channel,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", gvk.Kind)
			os.Exit(1)
		}

	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
