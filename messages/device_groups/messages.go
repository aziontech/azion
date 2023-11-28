package device_groups

var (
	// [ device groups ]
	DeviceGroupsUsage            = "device_groups"
	DeviceGroupsShortDescription = "Device groups is an Edge Application capability that allows you to identify the devices sending requests to your application."
	DeviceGroupsLongDescription  = "Device groups is an Edge Application capability that allows you to identify the devices sending requests to your application and categorize them into groups."
	DeviceGroupsFlagHelp         = "Displays more information about the Device Groups command"

	// [ list ]
	DeviceGroupsListUsage                 = "list [flags]"
	DeviceGroupsListShortDescription      = "Displays your device groups"
	DeviceGroupsListLongDescription       = "Displays all device groups related to a specific Edge Application"
	DeviceGroupsListHelpFlag              = "Displays more information about the list subcommand"
	DeviceGroupsListFlagEdgeApplicationID = "Unique identifier for an Edge Application."

	// [ delete ]
	DeviceGroupsDeleteUsage            = "delete [flags]"
	DeviceGroupsDeleteShortDescription = "Deletes a device group"
	DeviceGroupsDeleteLongDescription  = "Deletes a device group based on the given '--group-id' and '--application-id'"
	DeviceGroupsDeleteOutputSuccess    = "Device group %d was successfully deleted\n"
	DeviceGroupsDeleteHelpFlag         = "Displays more information about the delete subcommand"

	// describe cmd
	DeviceGroupsDescribeUsage            = "describe --application-id <application_id> --group-id <group_id> [flags]"
	DeviceGroupsDescribeShortDescription = "Returns the information related to a specific device group"
	DeviceGroupsDescribeLongDescription  = "Returns the information related to a specific device group, informed through the flag '--group-id' in detail"
	DeviceGroupsDescribeFlagOut          = "Exports the output of the subcommand 'describe' to the given file path <file_path/file_name.ext>"
	DeviceGroupsDescribeFlagFormat       = "Changes the output format passing the json value to the flag. Example '--format json'"
	DeviceGroupsDescribeHelpFlag         = "Displays more information about the describe subcommand"
	DeviceGroupsFileWritten              = "File successfully written to: %s\n"

	//update command
	DeviceGroupsUpdateUsage            = "update [flags]"
	DeviceGroupsUpdateShortDescription = "Updates a device group"
	DeviceGroupsUpdateLongDescription  = "Updates a device group based on given attributes to be used in Edge Applications"
	DeviceGroupsUpdateFlagName         = "The device group name"
	DeviceGroupsUpdateFlagUserAgent    = "The device group flag user agent"
	DeviceGroupsUpdateFlagIn           = "Path to a JSON file containing the attributes of the  device group that will be created; you can use - for reading from stdin"
	DeviceGroupsUpdateOutputSuccess    = "Device Group %d was updated\n"

	// [ create ]
	DeviceGroupsCreateUsage                 = "create [flags]"
	DeviceGroupsCreateShortDescription      = "Creates a new device group"
	DeviceGroupsCreateLongDescription       = "Creates a device group based on given attributes to be used in an Edge Application"
	DeviceGroupsCreateFlagEdgeApplicationId = "Unique identifier for an Edge Application"
	DeviceGroupsCreateFlagName              = "The name of your device group"
	DeviceGroupsCreateFlagUserAgent         = "The regex to match against the User-Agent header"
	DeviceGroupsCreateFlagIn                = "Path to a JSON file containing the attributes of the device group that will be created; you can use - for reading from stdin"
	DeviceGroupsCreateOutputSuccess         = "Created device group with ID %d\n"
	DeviceGroupsCreateHelpFlag              = "Displays more information about the create subcommand"

	ApplicationFlagId = "Unique identifier for the Edge Application that implements this device group. The '--application-id' flag is required"
	DeviceGroupFlagId = "Unique identifier for a device group. The '--group-id' flag is required"
)
