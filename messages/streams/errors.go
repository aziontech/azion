package streams

import "errors"

var (
	ErrorGetStream          = errors.New("Failed to get the stream: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetStreams         = errors.New("Failed to list the streams: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorFailToDeleteStream = errors.New("Failed to delete the stream: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCreateStream       = errors.New("Failed to create stream: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorUpdateStream       = errors.New("Failed to update the stream: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorConvertStreamId    = errors.New("Invalid --stream-id flag provided. The value must be an integer. Run the command 'azion describe streams --help' to display more information and try again")
)
