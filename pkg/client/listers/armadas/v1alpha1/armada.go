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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/chmouel/armadas/pkg/apis/armadas/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ArmadaLister helps list Armadas.
// All objects returned here must be treated as read-only.
type ArmadaLister interface {
	// List lists all Armadas in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Armada, err error)
	// Armadas returns an object that can list and get Armadas.
	Armadas(namespace string) ArmadaNamespaceLister
	ArmadaListerExpansion
}

// armadaLister implements the ArmadaLister interface.
type armadaLister struct {
	indexer cache.Indexer
}

// NewArmadaLister returns a new ArmadaLister.
func NewArmadaLister(indexer cache.Indexer) ArmadaLister {
	return &armadaLister{indexer: indexer}
}

// List lists all Armadas in the indexer.
func (s *armadaLister) List(selector labels.Selector) (ret []*v1alpha1.Armada, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Armada))
	})
	return ret, err
}

// Armadas returns an object that can list and get Armadas.
func (s *armadaLister) Armadas(namespace string) ArmadaNamespaceLister {
	return armadaNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ArmadaNamespaceLister helps list and get Armadas.
// All objects returned here must be treated as read-only.
type ArmadaNamespaceLister interface {
	// List lists all Armadas in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Armada, err error)
	// Get retrieves the Armada from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Armada, error)
	ArmadaNamespaceListerExpansion
}

// armadaNamespaceLister implements the ArmadaNamespaceLister
// interface.
type armadaNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Armadas in the indexer for a given namespace.
func (s armadaNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Armada, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Armada))
	})
	return ret, err
}

// Get retrieves the Armada from the indexer for a given namespace and name.
func (s armadaNamespaceLister) Get(name string) (*v1alpha1.Armada, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("armada"), name)
	}
	return obj.(*v1alpha1.Armada), nil
}