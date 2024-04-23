package autorest

type goTestOption struct {
	MockTest     string
	Example      string
	ScenarioTest string
	Sample       string
	FakeTest     string
}

var GOTestOption = goTestOption{
	MockTest:     "--testmodeler.generate-mock-test",
	Example:      "--testmodeler.generate-sdk-example",
	ScenarioTest: "--testmodeler.generate-scenario-test",
	Sample:       "--testmodeler.generate-sdk-sample",
	FakeTest:     "--testmodeler.generate-fake-test",
}
