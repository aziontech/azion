package edgeconnector

import "errors"

var (
	ErrorGetConnector          = errors.New("Failed to get the Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetConnectors         = errors.New("Failed to list the Connectors: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorFailToDeleteConnector = errors.New("Failed to delete the Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateConnector       = errors.New("Failed to create Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateConnector       = errors.New("Failed to update the Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
)
