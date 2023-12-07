/*
Copyright 2020 The Knative Authors

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

package reconciler

import (
	"context"

	"k8s.io/client-go/kubernetes"

	armadav1alpha1 "github.com/chmouel/armadas/pkg/apis/armadas/v1alpha1"
	fireconciler "github.com/chmouel/armadas/pkg/client/injection/reconciler/armadas/v1alpha1/fire"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
)

// Reconciler implements simpledeploymentreconciler.Interface for
// SimpleDeployment resources.
type Reconciler struct {
	kubeclient kubernetes.Interface
}

// Check that our Reconciler implements Interface
var _ fireconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, d *armadav1alpha1.Fire) reconciler.Event {
	// This logger has all the context necessary to identify which resource is being reconciled.
	logger := logging.FromContext(ctx)
	logger.Infof("Let's do a reconcilation  my friend: %v", d)
	return nil
}
