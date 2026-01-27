package scrapper

// PlatformStats represents statistics for a specific platform
type PlatformStats struct {
	Rating             int
	MaxRating          int
	ProblemsSolved     int
	EasyProblemsSolved int
	MedProblemsSolved  int
	HardProblemsSolved int
	ContestsAttended   int
	GlobalRank         int
}
