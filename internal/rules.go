package internal

type Rule int

const (
	RuleTrailingNewlinesList = 1 << iota
	RuleTrailingNewlinesMap
)
