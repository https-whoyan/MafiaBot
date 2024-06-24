package models

type TestRenameUserProvider struct{}

func (rP *TestRenameUserProvider) RenameUser(channelIID string, userServerID string, newNick string) error {
	return nil
}

var TestRenameUserProviderInstance = &TestRenameUserProvider{}
