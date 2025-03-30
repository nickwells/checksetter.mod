package checksetter

const (
	errFmtUnknownFunc  = "unknown function: %q"
	errFmtUnknownMaker = "unknown maker: %q"
	errFmtBadInt       = "couldn't make an int from %q: %w"
	errFmtBadFloat     = "couldn't make a float from %q: %w"
	errFmtNotABasicLit = "the expression isn't a BasicLit, it's a %T"
)
