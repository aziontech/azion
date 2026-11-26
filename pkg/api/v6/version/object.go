package version

// Object is a versioned resource as the v6 clients return it, in the way the
// v4 clients return a small response interface. Every SDK model of a
// versioned resource has these getters.
type Object interface {
	GetId() int64
	GetName() string
	GetVersionId() string
	GetVersionState() string
}

// Item is a sub-resource inside a version (a rule, a cache setting, a
// function instance, a WAF exception) as the v6 clients return it.
type Item interface {
	GetId() int64
	GetName() string
}

// Meta is one version of a resource as the version endpoints return it. Since
// SDK 0.3.0 those return their own model (ApplicationVersion and so on),
// which carries the version's id and state but not the resource's id.
type Meta interface {
	GetVersionId() string
	GetVersionState() string
}
