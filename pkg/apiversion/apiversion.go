// Package apiversion names the generations of the Azion API that the CLI can
// serve, and resolves which one an account is entitled to.
//
// The generation is resolved once at startup and carried on cmdutil.Factory, so
// the rest of the CLI reads it instead of re-deriving it from account flags.
//
// This package must not import pkg/token: token depends on cmdutil, which
// depends on this package. Persistence of a resolved version is expressed here
// as the dependency-free Cache type and mapped to token.Settings by callers.
package apiversion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Version identifies a generation of the Azion API.
//
// A new generation is added as a new constant here plus a new entry in
// generations, never as another boolean branch.
type Version string

const (
	// V3 is the legacy generation, served by pkg/v3commands and pkg/v3api.
	V3 Version = "v3"
	// V4 is the current generation, served by pkg/cmd and pkg/api.
	V4 Version = "v4"
)

// String returns the wire form of the version ("v3", "v4"). This is the value
// persisted in profile settings and reported in metrics, so it must stay stable.
func (v Version) String() string {
	return string(v)
}

// Client flags that mark an account as being on API v4. Either one is
// sufficient; accounts carrying neither stay on the legacy generation.
const (
	BlockAPIV4IncompatibleEndpoints = "block_apiv4_incompatible_endpoints"
	BlockAPIV3Access                = "block_apiv3_access"
)

// V4Flags is the set of client flags that entitle an account to API v4.
var V4Flags = []string{BlockAPIV4IncompatibleEndpoints, BlockAPIV3Access}

// AccountInfo is the subset of the SSO account payload the CLI reads.
type AccountInfo struct {
	ClientFlags []string `json:"client_flags"`
}

// HasAnyFlag reports whether the account carries at least one of the given
// client flags.
func (a AccountInfo) HasAnyFlag(names ...string) bool {
	for _, flag := range a.ClientFlags {
		for _, name := range names {
			if flag == name {
				return true
			}
		}
	}
	return false
}

// generation pairs a version with the test for whether an account may use it.
type generation struct {
	version  Version
	entitled func(AccountInfo) bool
}

// generations is the resolution order, most recent generation first. An account
// is served the first generation it is entitled to, so adding a generation means
// adding an entry here rather than another boolean branch.
//
// Entitlement to v4 is opt-in: the account must carry one of V4Flags. The legacy
// entry is unconditional, so it is also the answer for an account with no flags.
var generations = []generation{
	{
		version:  V4,
		entitled: func(info AccountInfo) bool { return info.HasAnyFlag(V4Flags...) },
	},
	{
		version:  V3,
		entitled: func(AccountInfo) bool { return true },
	},
}

// Select walks the resolution order and returns the newest generation the
// account is entitled to. The last entry is unconditional, so this always
// returns a version.
func Select(info AccountInfo) Version {
	for _, gen := range generations {
		if gen.entitled(info) {
			return gen.version
		}
	}
	return V3
}

// Resolve asks the SSO service which generation the given credential's account
// is on. Callers decide what an empty token or a failed lookup should mean.
func Resolve(client *http.Client, authURL, token string) (Version, error) {
	info, err := FetchAccountInfo(client, authURL, token)
	if err != nil {
		return "", err
	}
	return Select(info), nil
}

// FetchAccountInfo reads the account's client flags from the SSO service.
func FetchAccountInfo(client *http.Client, authURL, token string) (AccountInfo, error) {
	var info AccountInfo

	req, err := http.NewRequest("GET", authURL+"/account/info", nil)
	if err != nil {
		return info, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json; version=1")
	req.Header.Set("Authorization", "Token "+token)

	resp, err := client.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return info, fmt.Errorf("non-200 response: %d, body: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return info, err
	}

	return info, nil
}
