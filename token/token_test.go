package token

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/mocks/configure"
)

func init() {
	os.Setenv("SSO_MODE", "development")
}

func Test_Authenticate(t *testing.T) {

	mockReturn := configure.MockClient{}
	u, _ := url.Parse("http://localhost/apitoken?token=tokenTeste")
	request := http.Request{
		URL: u,
	}

	mockReturn.Get(&request)

	fmt.Println("Testing...")
}
