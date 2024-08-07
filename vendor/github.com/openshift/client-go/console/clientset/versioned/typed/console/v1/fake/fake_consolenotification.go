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

// FakeConsoleNotifications implements ConsoleNotificationInterface
type FakeConsoleNotifications struct {
	Fake *FakeConsoleV1
}

var consolenotificationsResource = schema.GroupVersionResource{Group: "console.openshift.io", Version: "v1", Resource: "consolenotifications"}

var consolenotificationsKind = schema.GroupVersionKind{Group: "console.openshift.io", Version: "v1", Kind: "ConsoleNotification"}

// Get takes name of the consoleNotification, and returns the corresponding consoleNotification object, and an error if there is any.
func (c *FakeConsoleNotifications) Get(ctx context.Context, name string, options v1.GetOptions) (result *consolev1.ConsoleNotification, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(consolenotificationsResource, name), &consolev1.ConsoleNotification{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleNotification), err
}

// List takes label and field selectors, and returns the list of ConsoleNotifications that match those selectors.
func (c *FakeConsoleNotifications) List(ctx context.Context, opts v1.ListOptions) (result *consolev1.ConsoleNotificationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(consolenotificationsResource, consolenotificationsKind, opts), &consolev1.ConsoleNotificationList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &consolev1.ConsoleNotificationList{ListMeta: obj.(*consolev1.ConsoleNotificationList).ListMeta}
	for _, item := range obj.(*consolev1.ConsoleNotificationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested consoleNotifications.
func (c *FakeConsoleNotifications) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(consolenotificationsResource, opts))
}

// Create takes the representation of a consoleNotification and creates it.  Returns the server's representation of the consoleNotification, and an error, if there is any.
func (c *FakeConsoleNotifications) Create(ctx context.Context, consoleNotification *consolev1.ConsoleNotification, opts v1.CreateOptions) (result *consolev1.ConsoleNotification, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(consolenotificationsResource, consoleNotification), &consolev1.ConsoleNotification{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleNotification), err
}

// Update takes the representation of a consoleNotification and updates it. Returns the server's representation of the consoleNotification, and an error, if there is any.
func (c *FakeConsoleNotifications) Update(ctx context.Context, consoleNotification *consolev1.ConsoleNotification, opts v1.UpdateOptions) (result *consolev1.ConsoleNotification, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(consolenotificationsResource, consoleNotification), &consolev1.ConsoleNotification{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleNotification), err
}

// Delete takes name of the consoleNotification and deletes it. Returns an error if one occurs.
func (c *FakeConsoleNotifications) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(consolenotificationsResource, name, opts), &consolev1.ConsoleNotification{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeConsoleNotifications) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(consolenotificationsResource, listOpts)

	_, err := c.Fake.Invokes(action, &consolev1.ConsoleNotificationList{})
	return err
}

// Patch applies the patch and returns the patched consoleNotification.
func (c *FakeConsoleNotifications) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *consolev1.ConsoleNotification, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(consolenotificationsResource, name, pt, data, subresources...), &consolev1.ConsoleNotification{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleNotification), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied consoleNotification.
func (c *FakeConsoleNotifications) Apply(ctx context.Context, consoleNotification *applyconfigurationsconsolev1.ConsoleNotificationApplyConfiguration, opts v1.ApplyOptions) (result *consolev1.ConsoleNotification, err error) {
	if consoleNotification == nil {
		return nil, fmt.Errorf("consoleNotification provided to Apply must not be nil")
	}
	data, err := json.Marshal(consoleNotification)
	if err != nil {
		return nil, err
	}
	name := consoleNotification.Name
	if name == nil {
		return nil, fmt.Errorf("consoleNotification.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(consolenotificationsResource, *name, types.ApplyPatchType, data), &consolev1.ConsoleNotification{})
	if obj == nil {
		return nil, err
	}
	return obj.(*consolev1.ConsoleNotification), err
}
