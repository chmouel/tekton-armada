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

	jobinformer "github.com/chmouel/armadas/pkg/client/injection/informers/armadas/v1alpha1/job"
	jobreconciler "github.com/chmouel/armadas/pkg/client/injection/reconciler/armadas/v1alpha1/job"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(ctx context.Context, _ configmap.Watcher) *controller.Impl {
	jobInformer := jobinformer.Get(ctx)

	r := &Reconciler{
		kubeclient: kubeclient.Get(ctx),
	}
	impl := jobreconciler.NewImpl(ctx, r)

	_, _ = jobInformer.Informer().AddEventHandler(
		controller.HandleAll(impl.Enqueue),
	)

	return impl
}
