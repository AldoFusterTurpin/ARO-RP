package openshiftcluster

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/Azure/go-autorest/autorest/azure"

	"github.com/jim-minter/rp/pkg/util/subnet"
)

func (m *Manager) Delete(ctx context.Context) error {
	r, err := azure.ParseResourceID(m.oc.ID)
	if err != nil {
		return err
	}

	m.log.Printf("deleting dns")
	_, err = m.recordsets.Delete(ctx, os.Getenv("RESOURCEGROUP"), m.domain, "api."+m.oc.Properties.DomainName, dns.CNAME, "")
	if err != nil {
		return err
	}

	// TODO: ideally we would do this after all the VMs have been deleted
	for subnetID, nsgID := range map[string]string{
		m.oc.Properties.MasterProfile.SubnetID:     "/subscriptions/" + r.SubscriptionID + "/resourceGroups/" + m.oc.Properties.ResourceGroup + "/providers/Microsoft.Network/networkSecurityGroups/" + m.oc.Properties.InfraID + "-controlplane-nsg",
		m.oc.Properties.WorkerProfiles[0].SubnetID: "/subscriptions/" + r.SubscriptionID + "/resourceGroups/" + m.oc.Properties.ResourceGroup + "/providers/Microsoft.Network/networkSecurityGroups/" + m.oc.Properties.InfraID + "-node-nsg",
	} {
		// TODO: there is probably an undesirable race condition here - check if etags can help.
		s, err := subnet.Get(ctx, &m.oc.Properties.ServicePrincipalProfile, subnetID)
		if err != nil {
			return err
		}

		if s.SubnetPropertiesFormat == nil ||
			s.SubnetPropertiesFormat.NetworkSecurityGroup == nil ||
			!strings.EqualFold(*s.SubnetPropertiesFormat.NetworkSecurityGroup.ID, nsgID) {
			continue
		}

		s.SubnetPropertiesFormat.NetworkSecurityGroup = nil

		m.log.Printf("removing network security group from subnet %s", subnetID)
		err = subnet.CreateOrUpdate(ctx, &m.oc.Properties.ServicePrincipalProfile, subnetID, s)
		if err != nil {
			return err
		}
	}

	resp, err := m.groups.CheckExistence(ctx, m.oc.Properties.ResourceGroup)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return nil
	}

	m.log.Printf("deleting resource group %s", m.oc.Properties.ResourceGroup)
	future, err := m.groups.Delete(ctx, m.oc.Properties.ResourceGroup)
	if err != nil {
		return err
	}

	m.log.Print("waiting for resource group deletion")
	return future.WaitForCompletionRef(ctx, m.groups.Client)
}
