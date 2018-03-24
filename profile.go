package bot

import "fmt"

type ProfileField int

const (
	Username ProfileField = iota
	FirstName
	LastName
)

func (self Bot) UpdateProfileField(profileField ProfileField, newValue string) {
	switch profileField {
	case Username:
		self.Server.Account.Username = newValue
	case FirstName:
		self.Server.Account.FirstName = newValue
	case LastName:
		self.Server.Account.LastName = newValue
	}
	self.Server.Account = self.Server.UpdateAccount(self.Server.Account)
}

func (self Bot) UpdateServerProfile() {
	fmt.Println("[SERVER] Checking if server profile is different from local configuration.")
	if self.Server.Account.FirstName != self.FirstName || self.Server.Account.LastName != self.LastName || self.Server.Account.Username != self.Username {
		fmt.Println("[SERVER] Server profile is different than local configuration, updating server profile.")
		self.Server.Account.FirstName = self.FirstName
		self.Server.Account.LastName = self.LastName
		self.Server.Account.Username = self.Username
		fmt.Println("[SERVER] New Server Value For First name:", self.Server.Account.FirstName)
		fmt.Println("[SERVER] New Server Value For Last name:", self.Server.Account.LastName)
		fmt.Println("[SERVER] New Server Value For Username:", self.Server.Account.Username)
		self.Server.Account = self.Server.UpdateAccount(self.Server.Account)
	} else {
		fmt.Println("[CONFIG] Values are the same, no need to update the server profile.")
	}
}
