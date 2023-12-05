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

package server

import (
	"context"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	armadainformer "github.com/chmouel/armadas/pkg/client/injection/informers/armadas/v1alpha1/armada"
	armadareconciler "github.com/chmouel/armadas/pkg/client/injection/reconciler/armadas/v1alpha1/armada"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	// Obtain an informer to both the main and child resources. These will be started by
	// the injection framework automatically. They'll keep a cached representation of the
	// cluster's state of the respective resource at all times.
	armadaInformer := armadainformer.Get(ctx)

	r := &Reconciler{
		// The client will be needed to create/delete Pods via the API.
		kubeclient: kubeclient.Get(ctx),
	}
	impl := armadareconciler.NewImpl(ctx, r)

	// Listen for events on the main resource and enqueue themselves.
	armadaInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	return impl
}
