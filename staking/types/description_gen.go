package types

import (
	"github.com/ethereum/go-ethereum/core/storage"
)

type DescriptionIdentitySetStorage struct {
	storage *storage.Storage
}

func (d *DescriptionIdentitySetStorage) ToDescriptionIdentitySet() DescriptionIdentitySet {
	return nil
}

func (d *DescriptionIdentitySetStorage) Size() int {
	return 0
}

func (d *DescriptionIdentitySetStorage) Contain(identity string) bool {
	return false
}

func (d *DescriptionIdentitySetStorage) Put(identity string) bool {
	return false
}

func (d *DescriptionIdentitySetStorage) Remove(identity string) {
}

func (d *DescriptionIdentitySetStorage) Keys() []string {
	return nil
}

type DescriptionStorage struct {
	storage *storage.Storage
}

func (d *DescriptionStorage) ToDescription() *Description {
	return nil
}

func (d *DescriptionStorage) GetName() string {
	return ""
}

func (d *DescriptionStorage) SetName(name string) {
}

func (d *DescriptionStorage) GetIdentity() string {
	return ""
}

func (d *DescriptionStorage) SetIdentity(identity string) {
}

func (d *DescriptionStorage) GetWebsite() string {
	return ""
}

func (d *DescriptionStorage) SetWebsite(website string) {
}

func (d *DescriptionStorage) GetSecurityContact() string {
	return ""
}

func (d *DescriptionStorage) SetSecurityContact(securityContact string) {
}
func (d *DescriptionStorage) GetDetails() string {
	return ""
}

func (d *DescriptionStorage) SetDetails(details string) {
}
