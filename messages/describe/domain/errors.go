package domain

import "errors"

var ErrorGetDomain = errors.New("Failed to describe the Domain: %s. Check your settings and try again. If the error persists, contact Azion support.")
