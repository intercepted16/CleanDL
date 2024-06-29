package main

type regexInfo struct {
	AgeThreshold int
	Destination  string
	DeleteFlag   bool
}

type regexPatternsJSON struct {
	Patterns regexPatterns
}

type regexPatterns map[string]regexInfo

type flagPointers struct {
	Pattern      *string
	AgeThreshold *int
	Destination  *string
	DeleteFlag   *bool
}
