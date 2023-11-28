package edgeapplication

import "errors"

var (
	ErrorUpdateApplication           = errors.New("Failed to update the Edge Application: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorApplicationAccelerationFlag = errors.New("Invalid --application-acceleration flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application  --help' to display more information and try again")
	ErrorCachingFlag                 = errors.New("Invalid --caching flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application --help' to display more information and try again")
	ErrorDeviceDetectionFlag         = errors.New("Invalid --device-detection flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application --help' to display more information and try again")
	ErrorEdgeFirewallFlag            = errors.New("Invalid --edge-firewall flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application --help' to display more information and try again")
	ErrorEdgeFunctionsFlag           = errors.New("Invalid --edge-functions flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application --help' to display more information and try again")
	ErrorImageOptimizationFlag       = errors.New("Invalid --image-optimization flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application --help' to display more information and try again")
	ErrorL2CachingFlag               = errors.New("Invalid --l2-caching flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application --help' to display more information and try again")
	ErrorLoadBalancerFlag            = errors.New("Invalid --load-balancer flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application --help' to display more information and try again")
	ErrorRawLogsFlag                 = errors.New("Invalid --raw-logs flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application --help' to display more information and try again")
	ErrorWebApplicationFirewallFlag  = errors.New("Invalid --webapp-firewall flag provided. The flag must have  'true' or 'false' values. Run the command 'azion update edge-application --help' to display more information and try again")
	ErrorNoFieldInformed             = errors.New("Inform at least one field to be updated. It's not possible to update an Edge Application without specifying the fields that will be updated. Run ‘azion update edge-application --help’ to display more information and try again.  If the error persists, contact Azion support.")
	ErrorConvertIdApplication        = errors.New("The application ID you provided is invalid. The value must be an integer. You may run the 'azion list edge-application' command to check your application ID")
)
