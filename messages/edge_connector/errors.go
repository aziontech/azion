package edgeconnector

import "errors"

var (
	ErrorGetConnector          = errors.New("Failed to get the Edge Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetConnectors         = errors.New("Failed to list the Edge Connectors: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorFailToDeleteConnector = errors.New("Failed to delete the Edge Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateConnector       = errors.New("Failed to create Edge Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateConnector       = errors.New("Failed to update the Edge Connector: %s. Check your settings and try again. If the error persists, contact Azion support")
)
