package device_groups

var (
	ApplicationFlagId = "Unique identifier for the edge application that implements this Device Group. The '--application-id' flag is required"
	DeviceGroupFlagId = "Unique identifier for a Device Group. The '--group-id' flag is required"

	// [ DeviceGroups ]
	DeviceGroupsUsage            = "device_groups"
	DeviceGroupsShortDescription = "Device Groups api"
	DeviceGroupsLongDescription  = "Device Groups api"
	DeviceGroupsFlagHelp         = "Displays more information about the Device Groups command"

	// [ delete ]
	DeviceGroupsDeleteUsage            = "delete [flags]"
	DeviceGroupsDeleteShortDescription = "Deletes a Device Group"
	DeviceGroupsDeleteLongDescription  = "Deletes a Device Group based on the given '--group-id' and '--application-id'"
	DeviceGroupsDeleteOutputSuccess    = "Device Group %d was successfully deleted\n"
	DeviceGroupsDeleteHelpFlag         = "Displays more information about the delete subcommand"
)
