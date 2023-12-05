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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/chmouel/armadas/pkg/apis/armadas/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFires implements FireInterface
type FakeFires struct {
	Fake *FakeGithubV1alpha1
	ns   string
}

var firesResource = v1alpha1.SchemeGroupVersion.WithResource("fires")

var firesKind = v1alpha1.SchemeGroupVersion.WithKind("Fire")

// Get takes name of the fire, and returns the corresponding fire object, and an error if there is any.
func (c *FakeFires) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Fire, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(firesResource, c.ns, name), &v1alpha1.Fire{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Fire), err
}

// List takes label and field selectors, and returns the list of Fires that match those selectors.
func (c *FakeFires) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FireList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(firesResource, firesKind, c.ns, opts), &v1alpha1.FireList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FireList{ListMeta: obj.(*v1alpha1.FireList).ListMeta}
	for _, item := range obj.(*v1alpha1.FireList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested fires.
func (c *FakeFires) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(firesResource, c.ns, opts))

}

// Create takes the representation of a fire and creates it.  Returns the server's representation of the fire, and an error, if there is any.
func (c *FakeFires) Create(ctx context.Context, fire *v1alpha1.Fire, opts v1.CreateOptions) (result *v1alpha1.Fire, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(firesResource, c.ns, fire), &v1alpha1.Fire{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Fire), err
}

// Update takes the representation of a fire and updates it. Returns the server's representation of the fire, and an error, if there is any.
func (c *FakeFires) Update(ctx context.Context, fire *v1alpha1.Fire, opts v1.UpdateOptions) (result *v1alpha1.Fire, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(firesResource, c.ns, fire), &v1alpha1.Fire{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Fire), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFires) UpdateStatus(ctx context.Context, fire *v1alpha1.Fire, opts v1.UpdateOptions) (*v1alpha1.Fire, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(firesResource, "status", c.ns, fire), &v1alpha1.Fire{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Fire), err
}

// Delete takes name of the fire and deletes it. Returns an error if one occurs.
func (c *FakeFires) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(firesResource, c.ns, name, opts), &v1alpha1.Fire{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFires) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(firesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FireList{})
	return err
}

// Patch applies the patch and returns the patched fire.
func (c *FakeFires) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Fire, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(firesResource, c.ns, name, pt, data, subresources...), &v1alpha1.Fire{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Fire), err
}
