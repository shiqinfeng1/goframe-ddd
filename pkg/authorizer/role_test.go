package authorizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoleFromString(t *testing.T) {
	cases := []struct {
		description string
		role        string
		expected    Role
	}{
		{
			description: "fails with invalid roles",
			role:        "N/A",
			expected:    RoleInvalid,
		},
		{
			description: "succeeds with 'owner'",
			role:        "owner",
			expected:    RoleOwner,
		},
		{
			description: "succeeds with 'administrator'",
			role:        "administrator",
			expected:    RoleAdministrator,
		},
		{
			description: "succeeds with 'operator'",
			role:        "operator",
			expected:    RoleOperator,
		},
		{
			description: "succeeds with 'observer'",
			role:        "observer",
			expected:    RoleObserver,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(tt *testing.T) {
			require.Equal(tt, tc.expected, RoleFromString(tc.role))
		})
	}
}

func TestRolePermissions(t *testing.T) {
	cases := []struct {
		description string
		role        Role
		expected    []Permission
	}{
		{
			description: "fails with invalid roles",
			role:        RoleInvalid,
			expected:    []Permission{},
		},
		{
			description: "succeeds with 'owner'",
			role:        RoleOwner,
			expected: []Permission{
				DeviceAccept,
				DeviceReject,
				DeviceRemove,
				DeviceConnect,
				DeviceRename,
				DeviceDetails,
				DeviceUpdate,
				DeviceCreateTag,
				DeviceUpdateTag,
				DeviceRemoveTag,
				DeviceRenameTag,
				DeviceDeleteTag,
				SessionPlay,
				SessionClose,
				SessionRemove,
				SessionDetails,
				FirewallCreate,
				FirewallEdit,
				FirewallRemove,
				FirewallAddTag,
				FirewallRemoveTag,
				FirewallUpdateTag,
				PublicKeyCreate,
				PublicKeyEdit,
				PublicKeyRemove,
				PublicKeyAddTag,
				PublicKeyRemoveTag,
				PublicKeyUpdateTag,
				NamespaceUpdate,
				NamespaceAddMember,
				NamespaceRemoveMember,
				NamespaceEditMember,
				NamespaceEnableSessionRecord,
				NamespaceDelete,
				BillingCreateCustomer,
				BillingChooseDevices,
				BillingAddPaymentMethod,
				BillingUpdatePaymentMethod,
				BillingRemovePaymentMethod,
				BillingCancelSubscription,
				BillingCreateSubscription,
				BillingGetSubscription,
				APIKeyCreate,
				APIKeyUpdate,
				APIKeyDelete,
				ConnectorDelete,
				ConnectorUpdate,
				ConnectorSet,
				TunnelsCreate,
				TunnelsDelete,
			},
		},
		{
			description: "succeeds with 'administrator'",
			role:        RoleAdministrator,
			expected: []Permission{
				DeviceAccept,
				DeviceReject,
				DeviceRemove,
				DeviceConnect,
				DeviceRename,
				DeviceDetails,
				DeviceUpdate,
				DeviceCreateTag,
				DeviceUpdateTag,
				DeviceRemoveTag,
				DeviceRenameTag,
				DeviceDeleteTag,
				SessionPlay,
				SessionClose,
				SessionRemove,
				SessionDetails,
				FirewallCreate,
				FirewallEdit,
				FirewallRemove,
				FirewallAddTag,
				FirewallRemoveTag,
				FirewallUpdateTag,
				PublicKeyCreate,
				PublicKeyEdit,
				PublicKeyRemove,
				PublicKeyAddTag,
				PublicKeyRemoveTag,
				PublicKeyUpdateTag,
				NamespaceUpdate,
				NamespaceAddMember,
				NamespaceRemoveMember,
				NamespaceEditMember,
				NamespaceEnableSessionRecord,
				APIKeyCreate,
				APIKeyUpdate,
				APIKeyDelete,
				ConnectorDelete,
				ConnectorUpdate,
				ConnectorSet,
				TunnelsCreate,
				TunnelsDelete,
			},
		},
		{
			description: "succeeds with 'operator'",
			role:        RoleOperator,
			expected: []Permission{
				DeviceAccept,
				DeviceReject,
				DeviceConnect,
				DeviceRename,
				DeviceDetails,
				DeviceUpdate,
				DeviceCreateTag,
				DeviceUpdateTag,
				DeviceRemoveTag,
				DeviceRenameTag,
				DeviceDeleteTag,
				SessionDetails,
			},
		},
		{
			description: "succeeds with 'observer'",
			role:        RoleObserver,
			expected: []Permission{
				DeviceConnect,
				DeviceDetails,
				SessionDetails,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(tt *testing.T) {
			require.Equal(tt, tc.expected, tc.role.Permissions())
		})
	}
}

func TestRolePreferences(t *testing.T) {
	cases := []struct {
		description string
		role        Role
		greater     []Role
		less        []Role
	}{
		{
			description: RoleInvalid.String(),
			role:        RoleInvalid,
			greater:     []Role{RoleOwner, RoleAdministrator, RoleOperator, RoleObserver},
			less:        []Role{},
		},
		{
			description: RoleOwner.String(),
			role:        RoleOwner,
			greater:     []Role{},
			less:        []Role{RoleAdministrator, RoleOperator, RoleObserver},
		},
		{
			description: RoleAdministrator.String(),
			role:        RoleAdministrator,
			greater:     []Role{RoleOwner},
			less:        []Role{RoleOperator, RoleObserver},
		},
		{
			description: RoleOperator.String(),
			role:        RoleOperator,
			greater:     []Role{RoleOwner, RoleAdministrator},
			less:        []Role{RoleObserver},
		},
		{
			description: RoleObserver.String(),
			role:        RoleObserver,
			greater:     []Role{RoleOwner, RoleAdministrator, RoleOperator},
			less:        []Role{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(tt *testing.T) {
			for _, r := range tc.greater {
				require.Equal(tt, false, tc.role.HasAuthority(r))
			}

			for _, r := range tc.less {
				require.Equal(tt, true, tc.role.HasAuthority(r))
			}
		})
	}
}
