package device_groups

// Used only by the v4 command tree.
var (
	// [ device groups ]
	DeviceGroupsUsage = "device-group"

	DeviceGroupsListShortDescription      = "Displays your device groups"
	DeviceGroupsListLongDescription       = "Displays all device groups related to a specific Application"
	DeviceGroupsListHelpFlag              = "Displays more information about the list subcommand"
	DeviceGroupsListFlagEdgeApplicationID = "Unique identifier for an Application."

	DeviceGroupsDeleteShortDescription = "Deletes a device group"
	DeviceGroupsDeleteLongDescription  = "Deletes a device group based on the given '--group-id' and '--application-id'"
	DeviceGroupsDeleteOutputSuccess    = "Device group %d was successfully deleted\n"
	DeviceGroupsDeleteHelpFlag         = "Displays more information about the delete subcommand"

	DeviceGroupsDescribeShortDescription = "Returns the information related to a specific device group"
	DeviceGroupsDescribeLongDescription  = "Returns the information related to a specific device group, informed through the flag '--group-id' in detail"
	DeviceGroupsDescribeHelpFlag         = "Displays more information about the describe subcommand"

	DeviceGroupsUpdateShortDescription = "Updates a device group"
	DeviceGroupsUpdateLongDescription  = "Updates a device group based on given attributes to be used in Applications"
	DeviceGroupsUpdateFlagName         = "The device group name"
	DeviceGroupsUpdateFlagUserAgent    = "The device group flag user agent"
	DeviceGroupsUpdateFlagIn           = "Path to a JSON file containing the attributes of the  device group that will be created; you can use - for reading from stdin"
	DeviceGroupsUpdateOutputSuccess    = "Device Group %d was updated\n"
	DeviceGroupsUpdateHelpFlag         = "Displays more information about the update subcommand"

	DeviceGroupsCreateShortDescription      = "Creates a new device group"
	DeviceGroupsCreateLongDescription       = "Creates a device group based on given attributes to be used in an Application"
	DeviceGroupsCreateFlagEdgeApplicationId = "Unique identifier for an Application"
	DeviceGroupsCreateFlagName              = "The name of your device group"
	DeviceGroupsCreateFlagUserAgent         = "The regex to match against the User-Agent header"
	DeviceGroupsCreateFlagIn                = "Path to a JSON file containing the attributes of the device group that will be created; you can use - for reading from stdin"
	DeviceGroupsCreateOutputSuccess         = "Created device group with ID %d\n"
	DeviceGroupsCreateHelpFlag              = "Displays more information about the create subcommand"

	ApplicationFlagId = "Unique identifier for the Application that implements this device group. The '--application-id' flag is required"
	DeviceGroupFlagId = "Unique identifier for a device group. The '--group-id' flag is required"

	// [ ask input prompts ]
	DeviceGroupsCreateAskInputApplicationID   = "Enter the ID of the Application the device group will be linked to:"
	DeviceGroupsCreateAskInputName            = "Enter the new device group's name:"
	DeviceGroupsCreateAskInputUserAgent       = "Enter the new device group's user agent:"
	DeviceGroupsListAskInputApplicationID     = "Enter the ID of the Application the device groups are linked to:"
	DeviceGroupsDescribeAskInputApplicationID = "Enter the ID of the Application the device group is linked to:"
	DeviceGroupsDescribeAskInputGroupID       = "Enter the ID of the device group you wish to describe:"
	DeviceGroupsUpdateAskInputApplicationID   = "Enter the ID of the Application the device group is linked to:"
	DeviceGroupsUpdateAskInputGroupID         = "Enter the ID of the device group you wish to update:"
	DeviceGroupsDeleteAskInputApplicationID   = "Enter the ID of the Application the device group is linked to:"
	DeviceGroupsDeleteAskInputGroupID         = "Enter the ID of the device group you wish to delete:"
)
