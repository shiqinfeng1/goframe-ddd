package authorizer

type Permission int

// 接口权限定义
const (
	DeviceAccept Permission = iota
	DeviceReject
	DeviceUpdate
	DeviceRemove
	DeviceConnect
	DeviceRename
	DeviceDetails
	DeviceCreateTag
	DeviceUpdateTag
	DeviceRemoveTag
	DeviceRenameTag
	DeviceDeleteTag

	SessionPlay
	SessionClose
	SessionRemove
	SessionDetails

	FirewallCreate
	FirewallEdit
	FirewallRemove
	FirewallAddTag
	FirewallRemoveTag
	FirewallUpdateTag

	PublicKeyCreate
	PublicKeyEdit
	PublicKeyRemove
	PublicKeyAddTag
	PublicKeyRemoveTag
	PublicKeyUpdateTag

	NamespaceUpdate
	NamespaceAddMember
	NamespaceRemoveMember
	NamespaceEditMember
	NamespaceEnableSessionRecord
	NamespaceDelete

	BillingCreateCustomer
	BillingChooseDevices
	BillingAddPaymentMethod
	BillingUpdatePaymentMethod
	BillingRemovePaymentMethod
	BillingCancelSubscription
	BillingCreateSubscription
	BillingGetPaymentMethod
	BillingGetSubscription

	APIKeyCreate
	APIKeyUpdate
	APIKeyDelete

	ConnectorDelete
	ConnectorUpdate
	ConnectorSet

	TunnelsCreate
	TunnelsDelete
)

// 每个角色具有的权限列表
var observerPermissions = []Permission{
	DeviceConnect,
	DeviceDetails,

	SessionDetails,
}

var operatorPermissions = []Permission{
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
}

var adminPermissions = []Permission{
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
}

var ownerPermissions = []Permission{
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
}
