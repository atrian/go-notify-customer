package handlers_test

type mockHandlerConfig struct {
}

func (m *mockHandlerConfig) GetTrustedSubnetAddress() string {
	return ""
}

func (m *mockHandlerConfig) GetDefaultResponseContentType() string {
	return "application/json"
}
