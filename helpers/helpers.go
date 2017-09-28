package helpers

// AssertNoError panics on error.
func AssertNoError(err error) {
	if err != nil {
		panic(err)
	}
}
