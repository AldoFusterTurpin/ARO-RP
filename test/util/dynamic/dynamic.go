package dynamic

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

// Client returns a client from a given unstructured object.
// It can be used when creating an object from a yaml file.
type Client interface {
	GetClient(obj *unstructured.Unstructured) (ResourceClient, error)
}

type client struct {
	dynamic.Interface
	mapping meta.RESTMapper
}

func NewDynamicClient(kubeConfig *rest.Config) (Client, error) {
	cli, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	cachedClient := memory.NewMemCacheClient(discoveryClient)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedClient)

	return &client{cli, mapper}, nil
}

func (d *client) getGVR(obj *unstructured.Unstructured) (*meta.RESTMapping, error) {
	var gvk = obj.GroupVersionKind()
	return d.mapping.RESTMapping(gvk.GroupKind(), gvk.Version)
}

// GetClient returns a non-namespaced or namespaced ResourceClient depending on a given object.
func (d *client) GetClient(obj *unstructured.Unstructured) (ResourceClient, error) {
	gvr, err := d.getGVR(obj)
	if err != nil {
		return nil, err
	}
	if obj.GetNamespace() != "" {
		return &resourceClient{d.Resource(gvr.Resource).Namespace(obj.GetNamespace())}, nil
	}
	return &resourceClient{d.Resource(gvr.Resource)}, nil
}

// ResourceClient is an interface that can be used for *K8sObjectWithRetry helper functions.
// The default dynamic ResourceInterface can't be used for *K8sObjectWithRetry
// because it has different method signatures from other clients interface.
// https://pkg.go.dev/k8s.io/client-go/dynamic#ResourceInterface
type ResourceClient interface {
	Get(ctx context.Context, name string, options metav1.GetOptions) (*unstructured.Unstructured, error)
	Create(ctx context.Context, obj *unstructured.Unstructured, options metav1.CreateOptions) (*unstructured.Unstructured, error)
	Delete(ctx context.Context, name string, options metav1.DeleteOptions) error
}

type resourceClient struct {
	client dynamic.ResourceInterface
}

func (d *resourceClient) Get(ctx context.Context, name string, options metav1.GetOptions) (*unstructured.Unstructured, error) {
	return d.client.Get(ctx, name, options)
}

func (d *resourceClient) Create(ctx context.Context, obj *unstructured.Unstructured, options metav1.CreateOptions) (*unstructured.Unstructured, error) {
	return d.client.Create(ctx, obj, options)
}

func (d *resourceClient) Delete(ctx context.Context, name string, options metav1.DeleteOptions) error {
	return d.client.Delete(ctx, name, options)
}
