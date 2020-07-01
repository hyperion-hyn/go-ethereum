package types

import (
	"github.com/pkg/errors"
)

type DescriptionIdentitySet map[string]struct{} // identity

// Description - some possible IRL connections
type Description struct {
	Name            string `json:"name"`             // name
	Identity        string `json:"identity"`         // optional identity signature (ex. UPort or Keybase)
	Website         string `json:"website"`          // optional website link
	SecurityContact string `json:"security-contact"` // optional security contact info
	Details         string `json:"details"`          // optional details
}

// EnsureLength ensures the length of a validator's description.
func (d *Description) EnsureLength() (*Description, error) {
	if len(d.Name) > MaxNameLength {
		return d, errors.Errorf(
			"exceed maximum name length %d %d", len(d.Name), MaxNameLength,
		)
	}
	if len(d.Identity) > MaxIdentityLength {
		return d, errors.Errorf(
			"exceed Maximum Length identity %d %d", len(d.Identity), MaxIdentityLength,
		)
	}
	if len(d.Website) > MaxWebsiteLength {
		return d, errors.Errorf(
			"exceed Maximum Length website %d %d", len(d.Website), MaxWebsiteLength,
		)
	}
	if len(d.SecurityContact) > MaxSecurityContactLength {
		return d, errors.Errorf(
			"exceed Maximum Length %d %d", len(d.SecurityContact), MaxSecurityContactLength,
		)
	}
	if len(d.Details) > MaxDetailsLength {
		return d, errors.Errorf(
			"exceed Maximum Length for details %d %d", len(d.Details), MaxDetailsLength,
		)
	}

	return d, nil
}

// UpdateDescription returns a new Description object with d1 as the base and the fields that's not empty in d2 updated
// accordingly. An error is returned if the resulting description fields have invalid length.
func UpdateDescription(curDecs *DescriptionStorage, newDesc *Description, identitySet *DescriptionIdentitySetStorage) error {
	if _, err := newDesc.EnsureLength(); err != nil {
		return err
	}

	if newDesc.Name != "" {
		curDecs.SetName(newDesc.Name)
	}
	if newDesc.Identity != "" {
		identitySet.Remove(curDecs.GetIdentity())
		curDecs.SetIdentity(newDesc.Identity)
		identitySet.Put(newDesc.Identity)
	}
	if newDesc.Website != "" {
		curDecs.SetWebsite(newDesc.Website)
	}
	if newDesc.SecurityContact != "" {
		curDecs.SetSecurityContact(newDesc.SecurityContact)
	}
	if newDesc.Details != "" {
		curDecs.SetDetails(newDesc.Details)
	}
	return nil
}
