// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	consolev1 "github.com/openshift/api/console/v1"
	applyconfigurationsconsolev1 "github.com/openshift/client-go/console/applyconfigurations/console/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeConsoleYAMLSamples implements ConsoleYAMLSampleInterface
type FakeConsoleYAMLSamples struct {
	Fake *FakeConsoleV1
}

var consoleyamlsamplesResource = schema.GroupVersionResource{Group: "console.openshift.io", Version: "v1", Resource: "consoleyamlsamples"}

var consoleyamlsamplesKind = schema.GroupVersionKind{Group: "console.openshift.io", Version: "v1", Kind: "ConsoleYAMLSample"}

// Get takes name of the consoleYAMLSample, and returns the corresponding consoleYAMLSample object, and an error if there is any.
func (c *FakeConsoleYAMLSamples) Get(ctx context.Context, name string, options v1.GetOptions) (result *consolev1.ConsoleYAMLSample, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(consoleyamlsamplesResource, name), &consolev1.ConsoleYAMLSample{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleYAMLSample), err
}

// List takes label and field selectors, and returns the list of ConsoleYAMLSamples that match those selectors.
func (c *FakeConsoleYAMLSamples) List(ctx context.Context, opts v1.ListOptions) (result *consolev1.ConsoleYAMLSampleList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(consoleyamlsamplesResource, consoleyamlsamplesKind, opts), &consolev1.ConsoleYAMLSampleList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &consolev1.ConsoleYAMLSampleList{ListMeta: obj.(*consolev1.ConsoleYAMLSampleList).ListMeta}
	for _, item := range obj.(*consolev1.ConsoleYAMLSampleList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested consoleYAMLSamples.
func (c *FakeConsoleYAMLSamples) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(consoleyamlsamplesResource, opts))
}

// Create takes the representation of a consoleYAMLSample and creates it.  Returns the server's representation of the consoleYAMLSample, and an error, if there is any.
func (c *FakeConsoleYAMLSamples) Create(ctx context.Context, consoleYAMLSample *consolev1.ConsoleYAMLSample, opts v1.CreateOptions) (result *consolev1.ConsoleYAMLSample, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(consoleyamlsamplesResource, consoleYAMLSample), &consolev1.ConsoleYAMLSample{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleYAMLSample), err
}

// Update takes the representation of a consoleYAMLSample and updates it. Returns the server's representation of the consoleYAMLSample, and an error, if there is any.
func (c *FakeConsoleYAMLSamples) Update(ctx context.Context, consoleYAMLSample *consolev1.ConsoleYAMLSample, opts v1.UpdateOptions) (result *consolev1.ConsoleYAMLSample, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(consoleyamlsamplesResource, consoleYAMLSample), &consolev1.ConsoleYAMLSample{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleYAMLSample), err
}

// Delete takes name of the consoleYAMLSample and deletes it. Returns an error if one occurs.
func (c *FakeConsoleYAMLSamples) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(consoleyamlsamplesResource, name, opts), &consolev1.ConsoleYAMLSample{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeConsoleYAMLSamples) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(consoleyamlsamplesResource, listOpts)

	_, err := c.Fake.Invokes(action, &consolev1.ConsoleYAMLSampleList{})
	return err
}

// Patch applies the patch and returns the patched consoleYAMLSample.
func (c *FakeConsoleYAMLSamples) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *consolev1.ConsoleYAMLSample, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(consoleyamlsamplesResource, name, pt, data, subresources...), &consolev1.ConsoleYAMLSample{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleYAMLSample), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied consoleYAMLSample.
func (c *FakeConsoleYAMLSamples) Apply(ctx context.Context, consoleYAMLSample *applyconfigurationsconsolev1.ConsoleYAMLSampleApplyConfiguration, opts v1.ApplyOptions) (result *consolev1.ConsoleYAMLSample, err error) {
	if consoleYAMLSample == nil {
		return nil, fmt.Errorf("consoleYAMLSample provided to Apply must not be nil")
	}
	data, err := json.Marshal(consoleYAMLSample)
	if err != nil {
		return nil, err
	}
	name := consoleYAMLSample.Name
	if name == nil {
		return nil, fmt.Errorf("consoleYAMLSample.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(consoleyamlsamplesResource, *name, types.ApplyPatchType, data), &consolev1.ConsoleYAMLSample{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleYAMLSample), err
}
