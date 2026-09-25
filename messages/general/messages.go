package general

// Used by both the v3 and the v4 command trees.
var (
	//List flags used by all azion apis
	ApiListFlagDetails  = "Displays all relevant fields when listing"
	ApiListFlagOrderBy  = "Sorts the output based on the selected field"
	ApiListFlagPage     = "Returns a page of the list according to its number"
	ApiListFlagPageSize = "Defines how many items should be returned per page"
	ApiListFlagNextPage = "token to next page"
	ApiListFlagFilter   = "Filters items by their name"
	CliVersion          = "Azion CLI %s"
)

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ApiListFlagSort = "Defines the order of the items on the list; options <asc|desc>"
)
