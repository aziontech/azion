package device_groups

var (
	// [ device groups ]
	DeviceGroupsUsage            = "device_groups"
	DeviceGroupsShortDescription = "Device groups is an Edge Application capability that allows you to identify the devices sending requests to your application."
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

	// describe cmd
	DeviceGroupsDescribeUsage            = "describe --application-id <application_id> --group-id <group_id> [flags]"
	DeviceGroupsDescribeShortDescription = "Returns the information related to the Device Group"
	DeviceGroupsDescribeLongDescription  = "Returns the information related to the Device Group, informed through the flag '--group-id' in detail"
	DeviceGroupsDescribeFlagOut          = "Exports the output of the subcommand 'describe' to the given file path <file_path/file_name.ext>"
	DeviceGroupsDescribeFlagFormat       = "Changes the output format passing the json value to the flag. Example '--format json'"
	DeviceGroupsDescribeHelpFlag         = "Displays more information about the describe subcommand"
	DeviceGroupsFileWritten              = "File successfully written to: %s\n"
  
	//update command
	DeviceGroupsUpdateUsage            = "update [flags]"
	DeviceGroupsUpdateShortDescription = "Updates a device group"
	DeviceGroupsUpdateLongDescription  = "Updates a device group based on given attributes to be used in edge applications"
	DeviceGroupsUpdateFlagName         = "The device group name"
	DeviceGroupsUpdateFlagUserAgent    = "The device group flag user agent"
	DeviceGroupsUpdateFlagIn           = "Path to a JSON file containing the attributes of the  device group that will be created; you can use - for reading from stdin"
	DeviceGroupsUpdateOutputSuccess    = "Device Group %d was updated updated\n"

	// [ create ]
	DeviceGroupsCreateUsage                 = "create [flags]"
	DeviceGroupsCreateShortDescription      = "Creates a new device groups"
	DeviceGroupsCreateLongDescription       = "Creates an device groups based on given attributes to be used in edge applications"
	DeviceGroupsCreateFlagEdgeApplicationId = "Unique identifier for an edge application"
	DeviceGroupsCreateFlagName              = "The device group name"
	DeviceGroupsCreateFlagUserAgent         = "the device group flag user agent"
	DeviceGroupsCreateFlagIn                = "Path to a JSON file containing the attributes of the  device group that will be created; you can use - for reading from stdin"
	DeviceGroupsCreateOutputSuccess         = "Created  device group with ID %d\n"
	DeviceGroupsCreateHelpFlag              = "Displays more information about the create subcommand"

	ApplicationFlagId = "Unique identifier for the edge application that implements this Device Group. The '--application-id' flag is required"
	DeviceGroupFlagId = "Unique identifier for a Device Group. The '--group-id' flag is required"
)
