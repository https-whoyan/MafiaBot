package models

type TestRenameUserProvider struct{}

var TestRenameUserProviderInstance = &TestRenameUserProvider{}

func (rP *TestRenameUserProvider) RenameUser(_ string, _ string, _ string) error { return nil }
