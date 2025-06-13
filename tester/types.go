package tester

type TestSuite struct {
	BaseURL   string     `yaml:"base_url"`
	TenantID  string     `yaml:"tenant_id"`
	ChannelID string     `yaml:"channel_id"`
	Secret    string     `yaml:"channel_secret"`
	TestCases []TestCase `yaml:"test_cases"`
}

type TestCase struct {
	Name           string `yaml:"name"`
	PayloadFile    string `yaml:"payload_file"`
	ExpectedStatus int    `yaml:"expected_status"`
}

type Result struct {
	Name     string
	Passed   bool
	Status   int
	Message  string
	Response string
}
