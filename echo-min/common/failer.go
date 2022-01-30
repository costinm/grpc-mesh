package common

// Failer is an interface to be provided to test functions of the form XXXOrFail. This is a
// substitute for testing.TB, which cannot be implemented outside of the testing
// package.
type Failer interface {
	Fail()
	FailNow()
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	TempDir() string
	Helper()
	Cleanup(func())
}

