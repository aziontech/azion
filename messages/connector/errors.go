package connector

import "errors"

var (
	ErrorGetConnector          = errors.New("Failed to get the Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetConnectors         = errors.New("Failed to list the Connectors: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorFailToDeleteConnector = errors.New("Failed to delete the Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateConnector       = errors.New("Failed to create Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateConnector       = errors.New("Failed to update the Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertConnectorId    = errors.New("Invalid --connector-id flag provided. The value must be an integer. Run the command 'azion create function-instance --help' to display more information and try again")
)
