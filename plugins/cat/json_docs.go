package plugins

// BuildJSON manages the json output from the cat build command results
type BuildJSON struct {
	Objects []struct {
		Build struct {
			Branch string
		}
		Result    string
		SbBuildID int
		Targets   string
	}
}

// CurrBuildJSON manages the json output from the cat current build results
type CurrBuildJSON struct {
	Changeset string
}

// RecommendedJSON manages the json output from the cat recommended command results
type RecommendedJSON struct {
	Objects []struct {
		CurrBuild string
	}
}

// TestrunJSON manages the json output from the cat testrun command results
type TestrunJSON struct {
	Objects []struct {
		ID         int
		Result     string
		ResultsDir string
	}
}
