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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Fire is spanish for yoplait
//
// +k8s:openapi-gen=true
type Fire struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec holds the desired state of the Armada (from the client).
	Spec FireSpec `json:"spec,omitempty"`

	// Status communicates the observed state of the Armada (from the controller).
	Status FireStatus `json:"status,omitempty"`
}

var (
	// Check that Armada can be validated and defaulted.
	_ apis.Validatable   = (*Fire)(nil)
	_ apis.Defaultable   = (*Fire)(nil)
	_ kmeta.OwnerRefable = (*Fire)(nil)
	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*Fire)(nil)
)

// FireSpec defines the desired state of the Fire, represented
// by a list of Fires.
type FireSpec struct {
	Tags  []string `json:"tags"`
	YAMLs []string `json:"yamls"`
}

// FireStatus is the status that makes it the best of the best
type FireStatus struct {
	duckv1.Status `json:",inline"`

	Accepted []*duckv1.SourceList `json:"address,omitempty"`
}

// FireList is a list of Armada resources
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FireList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Fire `json:"items"`
}

// GetStatus retrieves the status of the resource. Implements the KRShaped interface.
func (d *Fire) GetStatus() *duckv1.Status {
	return &d.Status.Status
}
