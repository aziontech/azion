package device_groups

var (
	// [ device groups ]
	DeviceGroupsUsage            = "device_groups"
	DeviceGroupsShortDescription = "Device groups is the original source of data."
	DeviceGroupsLongDescription  = "Device groups is the original source of data on edge platforms, where data is fetched when cache is not available."
	DeviceGroupsFlagHelp         = "Displays more information about the Device groups command"

	// [ list ]
	DeviceGroupsListUsage                 = "list [flags]"
	DeviceGroupsListShortDescription      = "Displays your device groups"
	DeviceGroupsListLongDescription       = "Displays all device groups related to your applications"
	DeviceGroupsListHelpFlag              = "Displays more information about the list subcommand"
	DeviceGroupsListFlagEdgeApplicationID = "Unique identifier for an edge application."

	// [ delete ]
	DeviceGroupsDeleteUsage            = "delete [flags]"
	DeviceGroupsDeleteShortDescription = "Deletes a Device Group"
	DeviceGroupsDeleteLongDescription  = "Deletes a Device Group based on the given '--group-id' and '--application-id'"
	DeviceGroupsDeleteOutputSuccess    = "Device Group %d was successfully deleted\n"
	DeviceGroupsDeleteHelpFlag         = "Displays more information about the delete subcommand"
  
  ApplicationFlagId = "Unique identifier for the edge application that implements this Device Group. The '--application-id' flag is required"
	DeviceGroupFlagId = "Unique identifier for a Device Group. The '--group-id' flag is required"
)
