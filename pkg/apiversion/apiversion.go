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
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// ErrUnauthorized reports that the authentication service rejected the
// credential. The account's generation is unknowable until the credential is
// fixed, and the command being run reports that problem in its own terms, so
// callers should stay quiet about the version rather than add noise.
var ErrUnauthorized = errors.New("credential rejected by the authentication service")

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

const (
	BlockAPIV4IncompatibleEndpoints = "block_apiv4_incompatible_endpoints"
	BlockAPIV3Access                = "block_apiv3_access"
)

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
// v4 is the default: an account only reaches the legacy entry by carrying the
// flag that blocks it from v4. The legacy entry is unconditional, so it is also
// the answer when no generation above it accepts the account.
var generations = []generation{
	{
		version: V4,
		entitled: func(info AccountInfo) bool {
			// Being blocked from v3 is the same statement from the other side,
			// and wins if both flags are somehow set: an account that cannot use
			// v3 at all is only ever worse off on the legacy tree.
			if info.HasAnyFlag(BlockAPIV3Access) {
				return true
			}
			return !info.HasAnyFlag(BlockAPIV4IncompatibleEndpoints)
		},
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
	spew.Dump(info.ClientFlags)
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
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return info, fmt.Errorf("%w: %d", ErrUnauthorized, resp.StatusCode)
		}
		return info, fmt.Errorf("non-200 response: %d, body: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return info, err
	}

	return info, nil
}
