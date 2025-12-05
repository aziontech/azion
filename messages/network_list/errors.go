package networklist

import "errors"

var (
	ErrorGetNetworkLists         = errors.New("Failed to list your network lists. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorFailToDeleteNetworkList = errors.New("Failed to delete Network List: %w")
	ErrorGetNetworkList          = errors.New("Failed to describe Network List: %w")
	ErrorCreateNetworkList       = errors.New("Failed to create Network List: %w")
	ErrorUpdateNetworkList       = errors.New("Failed to update Network List: %w")
	ErrorActiveFlag              = errors.New("Invalid value for --active flag")
	ErrorConvertNetworkListId    = errors.New("Invalid --network-list-id flag provided. The value must be an integer. Run the command 'azion delete network-list --help' to display more information and try again")
)
