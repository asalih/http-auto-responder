package utils

//["CONTAINS", "WILDCARD", "EXACT", "NOT", "REGEX"]
const (
	CONTAINS = "CONTAINS"
	WILDCARD = "WILDCARD"
	EXACT    = "EXACT"
	NOT      = "NOT"
	REGEX    = "REGEX"
)

//GetMatchType Gets the match type
func GetMatchType(mt string) string {
	if mt == NOT || mt == EXACT || mt == CONTAINS || mt == REGEX || mt == WILDCARD {
		return mt
	}
	return CONTAINS
}
