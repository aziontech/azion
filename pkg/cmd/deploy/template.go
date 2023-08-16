package deploy

import (
	"strings"
	"text/template"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

const jsCode = `
// Define the project type as a static site
self.__PROJECT_TYPE_PATTERN = "PROJECT_TYPE:STATIC_SITE";

// Listen for incoming events and respond with handleEvent function
addEventListener('fetch', event => {
  event.respondWith(handleEvent(event))
})

// Handle incoming events
async function handleEvent(event) {
  try {
    // Get the requested path from the event URL
    const request_path = new URL(event.request.url).pathname;

    // Get the version ID for the requested asset
    const version_id = "{{ .VersionId }}";

    /* Often web servers are configured to look for a default document when a directory is requested. 
    For example, if the server receives a request for http://example.com/directory/, it might 
    automatically look for a file named index.html or default.aspx within that directory, 
    and serve that file as the response.
    This behavior is configurable through the Edge Functions, and the default file names can vary.*/
    let asset_path;
    if (request_path === "/") {
      // If the requested path is just "/", construct the asset path with "/index.html"
      asset_path = version_id + "/index.html";
    } else if (request_path.endsWith("/")) {
      // If the requested path ends with a "/", concatenate the path with "index.html"
      asset_path = version_id + request_path + "index.html";
    } else {
      // For all other cases, use the requested path as the asset path
      asset_path = version_id + request_path;
    }

    // Construct the URL for the requested asset
    const asset_url = new URL(asset_path, "file://");
    // Return the fetch response for the asset
    return fetch(asset_url);

  } catch (e) {
    // If there is an error, return a Response object with the error message and a status code of 500
    return new Response(e.message || e.toString(), { status: 500 });
  }
}
`

func (cmd *DeployCmd) applyTemplate(conf *contracts.AzionApplicationOptions) (string, error) {
	tmpl, err := template.New("jsTemplate").Parse(jsCode)
	if err != nil {
		logger.Debug("Error while parsing template in javascript function", zap.Error(err))
		return "", utils.ErrorParsingModel
	}

	data := struct {
		VersionId string
	}{
		VersionId: conf.VersionID,
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		logger.Debug("Error while applying template to javascript function", zap.Error(err))
		return "", utils.ErrorExecTemplate
	}

	return result.String(), nil
}
