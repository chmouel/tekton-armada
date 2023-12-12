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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	httpClient *http.Client
}

func (r *Reconciler) fireToTheServer(ctx context.Context, spec armadav1alpha1.FireSpec, dest string) (*http.Response, error) {
	// encode spec to json
	jsonSpec, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(dest)
	if err != nil {
		return nil, err
	}
	// make a request
	request := &http.Request{
		Method: "POST",
		URL:    u,
		Body:   io.NopCloser(bytes.NewBuffer(jsonSpec)),
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}

	resp, err := r.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New("failed to parse response body")
		}
		return nil, fmt.Errorf("request rejected; status: %s; message: %s", resp.Status, respBody)
	}
	return resp, nil
}

// Check that our Reconciler implements Interface
var _ fireconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, d *armadav1alpha1.Fire) reconciler.Event {
	// This logger has all the context necessary to identify which resource is being reconciled.
	logger := logging.FromContext(ctx)
	logger.Infof("Let's do a reconcilation my friend: %v", d)
	_, err := r.fireToTheServer(ctx, d.Spec, "http//armada-server:3344")
	if err != nil {
		return err
	}
	return nil
}
