package cluster

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"testing"

	mgmtnetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/golang/mock/gomock"
	imageregistryv1 "github.com/openshift/api/imageregistry/v1"

	"github.com/Azure/ARO-RP/pkg/api"
	mock_subnet "github.com/Azure/ARO-RP/pkg/util/mocks/subnet"
)

var (
	subscriptionId    = "0000000-0000-0000-0000-000000000000"
	vnetResourceGroup = "vnet-rg"
	vnetName          = "vnet"
	subnetNameWorker  = "worker"
	subnetNameMaster  = "master"
)

func getValidSubnet(endpoints bool, state *mgmtnetwork.ProvisioningState) *mgmtnetwork.Subnet {
	s := &mgmtnetwork.Subnet{
		SubnetPropertiesFormat: &mgmtnetwork.SubnetPropertiesFormat{},
	}
	if endpoints {
		s.SubnetPropertiesFormat = &mgmtnetwork.SubnetPropertiesFormat{
			ServiceEndpoints: &[]mgmtnetwork.ServiceEndpointPropertiesFormat{
				{
					Service:   to.StringPtr("Microsoft.ContainerRegistry"),
					Locations: &[]string{"*"},
				},
				{
					Service:   to.StringPtr("Microsoft.Storage"),
					Locations: &[]string{"*"},
				},
			},
		}
		if state != nil {
			for i := range *s.SubnetPropertiesFormat.ServiceEndpoints {
				(*s.SubnetPropertiesFormat.ServiceEndpoints)[i].ProvisioningState = *state
			}
		}
	}
	return s
}

func TestEnableServiceEndpoints(t *testing.T) {
	ctx := context.Background()

	type test struct {
		name string
		oc   *api.OpenShiftCluster
		mock func(subnetMock *mock_subnet.MockManager, tt test)
	}

	for _, tt := range []test{
		{
			name: "nothing to do",
			oc: &api.OpenShiftCluster{
				Properties: api.OpenShiftClusterProperties{
					MasterProfile: api.MasterProfile{
						SubnetID: "/subscriptions/" + subscriptionId + "/resourceGroups/" + vnetResourceGroup + "/providers/Microsoft.Network/virtualNetworks/" + vnetName + "/subnets/" + subnetNameMaster,
					},
					WorkerProfiles: []api.WorkerProfile{
						{
							SubnetID: "/subscriptions/" + subscriptionId + "/resourceGroups/" + vnetResourceGroup + "/providers/Microsoft.Network/virtualNetworks/" + vnetName + "/subnets/" + subnetNameWorker,
						},
					},
				},
			},
			mock: func(subnetClient *mock_subnet.MockManager, tt test) {
				subnets := []string{
					tt.oc.Properties.MasterProfile.SubnetID,
				}
				for _, wp := range tt.oc.Properties.WorkerProfiles {
					subnets = append(subnets, wp.SubnetID)
				}

				for _, subnetId := range subnets {
					state := mgmtnetwork.Succeeded
					subnetClient.EXPECT().Get(gomock.Any(), subnetId).Return(getValidSubnet(true, &state), nil)
				}
			},
		},
		{
			name: "enable endpoints",
			oc: &api.OpenShiftCluster{
				Properties: api.OpenShiftClusterProperties{
					MasterProfile: api.MasterProfile{
						SubnetID: "/subscriptions/" + subscriptionId + "/resourceGroups/" + vnetResourceGroup + "/providers/Microsoft.Network/virtualNetworks/" + vnetName + "/subnets/" + subnetNameMaster,
					},
					WorkerProfiles: []api.WorkerProfile{
						{
							SubnetID: "/subscriptions/" + subscriptionId + "/resourceGroups/" + vnetResourceGroup + "/providers/Microsoft.Network/virtualNetworks/" + vnetName + "/subnets/" + subnetNameWorker,
						},
					},
				},
			},
			mock: func(subnetClient *mock_subnet.MockManager, tt test) {
				subnets := []string{
					tt.oc.Properties.MasterProfile.SubnetID,
				}
				for _, wp := range tt.oc.Properties.WorkerProfiles {
					subnets = append(subnets, wp.SubnetID)
				}

				for _, subnetId := range subnets {
					subnetClient.EXPECT().Get(gomock.Any(), subnetId).Return(getValidSubnet(false, nil), nil)
					subnetClient.EXPECT().CreateOrUpdate(gomock.Any(), subnetId, getValidSubnet(true, nil)).Return(nil)
				}
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			subnetClient := mock_subnet.NewMockManager(controller)

			tt.mock(subnetClient, tt)

			m := &manager{
				subnet: subnetClient,
				doc: &api.OpenShiftClusterDocument{
					OpenShiftCluster: tt.oc,
				},
			}

			// we don't test errors as all of them would be out of our control
			_ = m.enableServiceEndpoints(ctx)
		})
	}
}

func TestGetAccountName(t *testing.T) {
	type testFields struct {
		name                string
		registryConfig      *imageregistryv1.Config
		expectedAccountName string
		expectedError       string
	}

	testCases := []testFields{
		{
			name:                "should return empty string and error when registry config is nil",
			registryConfig:      nil,
			expectedAccountName: "",
			expectedError:       "image registry config is nil",
		},
		{
			name:                "should return empty string and error when registry config is empty",
			registryConfig:      &imageregistryv1.Config{},
			expectedAccountName: "",
			expectedError:       "azure storage field is nil in image registry config",
		},
		{
			name: "should return empty string and error when Spec field is empty",
			registryConfig: &imageregistryv1.Config{
				Spec: imageregistryv1.ImageRegistrySpec{},
			},
			expectedAccountName: "",
			expectedError:       "azure storage field is nil in image registry config",
		},
		{
			name: "should return the appropiate account name and no error when account name exists",
			registryConfig: &imageregistryv1.Config{
				Spec: imageregistryv1.ImageRegistrySpec{
					Storage: imageregistryv1.ImageRegistryConfigStorage{
						Azure: &imageregistryv1.ImageRegistryConfigStorageAzure{
							AccountName: "the_account_name",
						},
					},
				},
			},
			expectedAccountName: "the_account_name",
			expectedError:       "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accountName, err := getAccountName(tc.registryConfig)

			if tc.expectedAccountName != accountName {
				t.Fatalf("expected account name %v, but got %v", tc.expectedAccountName, accountName)
			}

			if err != nil && err.Error() != tc.expectedError || err == nil && tc.expectedError != "" {
				t.Fatalf("expected error %v, but got %v", tc.expectedError, err)
			}
		})
	}
}

func TestGetAccountNameMutator(t *testing.T) {
	type testFields struct {
		name                       string
		registryConfig             *imageregistryv1.Config
		doc                        *api.OpenShiftClusterDocument
		expectedAccountNameFromDoc string
		expectedError              string
	}

	testCases := []testFields{
		{
			name:                       "should return error when doc is nil",
			registryConfig:             nil,
			doc:                        nil,
			expectedAccountNameFromDoc: "",
			expectedError:              "OpenShift cluster document is nil",
		},
		{
			name:                       "should return error when doc.OpenShiftCluster is nil",
			registryConfig:             nil,
			doc:                        &api.OpenShiftClusterDocument{},
			expectedAccountNameFromDoc: "",
			expectedError:              "OpenShiftCluster info from OpenShift cluster document is nil",
		},
		{
			name:           "should propagate the error from getting the registry when that error is not nil",
			registryConfig: nil,
			doc: &api.OpenShiftClusterDocument{
				OpenShiftCluster: &api.OpenShiftCluster{},
			},
			expectedAccountNameFromDoc: "",
			expectedError:              "image registry config is nil",
		},
		{
			name: "should set account name of openShiftClusterDocument when account name of registry config has a valid value",
			registryConfig: &imageregistryv1.Config{
				Spec: imageregistryv1.ImageRegistrySpec{
					Storage: imageregistryv1.ImageRegistryConfigStorage{
						Azure: &imageregistryv1.ImageRegistryConfigStorageAzure{
							AccountName: "the_account_name",
						},
					},
				},
			},
			doc: &api.OpenShiftClusterDocument{
				OpenShiftCluster: &api.OpenShiftCluster{},
			},
			expectedAccountNameFromDoc: "the_account_name",
			expectedError:              "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accountNameMutator := getAccountNameMutator(context.Background(), tc.registryConfig)

			err := accountNameMutator(tc.doc)

			if err != nil && err.Error() != tc.expectedError ||
				err == nil && tc.expectedError != "" {
				t.Fatalf("expected error %v, but got %v", tc.expectedError, err)
			}

			if tc.expectedError == "" {
				docAccountName := tc.doc.OpenShiftCluster.Properties.ImageRegistryStorageAccountName

				if tc.expectedAccountNameFromDoc != docAccountName {
					t.Fatalf("expected account name from doc %v, but got %v", tc.expectedAccountNameFromDoc, docAccountName)
				}
			}
		})
	}
}
