package restaking

import "github.com/pkg/errors"

// UpdateDescription returns a new Description object with d1 as the base and the fields that's not empty in d2 updated
// accordingly. An error is returned if the resulting description fields have invalid length.
func (d *Description_) UpdateFrom(other Description_) error {
	if other.Name != "" {
		d.Name = other.Name
	}
	if other.Identity != "" {
		d.Identity = other.Identity
	}
	if other.Website != "" {
		d.Website = other.Website
	}
	if other.SecurityContact != "" {
		d.SecurityContact = other.SecurityContact
	}
	if other.Details != "" {
		d.Details = other.Details
	}
	return d.EnsureLength()
}

// EnsureLength ensures the length of a validator's description.
func (d *Description_) EnsureLength() error {
	if len(d.Name) > MaxNameLength {
		return errors.Errorf(
			"exceed maximum name length %d %d", len(d.Name), MaxNameLength,
		)
	}
	if len(d.Identity) > MaxIdentityLength {
		return errors.Errorf(
			"exceed Maximum Length identity %d %d", len(d.Identity), MaxIdentityLength,
		)
	}
	if len(d.Website) > MaxWebsiteLength {
		return errors.Errorf(
			"exceed Maximum Length website %d %d", len(d.Website), MaxWebsiteLength,
		)
	}
	if len(d.SecurityContact) > MaxSecurityContactLength {
		return errors.Errorf(
			"exceed Maximum Length %d %d", len(d.SecurityContact), MaxSecurityContactLength,
		)
	}
	if len(d.Details) > MaxDetailsLength {
		return errors.Errorf(
			"exceed Maximum Length for details %d %d", len(d.Details), MaxDetailsLength,
		)
	}
	return nil
}

// Storage_Description_
func (s *Storage_Description_) UpdateDescription(newDesc Description_) {
	if newDesc.Name != "" {
		s.Name().SetValue(newDesc.Name)
	}
	if newDesc.Identity != "" {
		s.Identity().SetValue(newDesc.Identity)
	}
	if newDesc.Website != "" {
		s.Website().SetValue(newDesc.Website)
	}
	if newDesc.SecurityContact != "" {
		s.SecurityContact().SetValue(newDesc.SecurityContact)
	}
	if newDesc.Details != "" {
		s.Details().SetValue(newDesc.Details)
	}
}

func (s *Storage_Description_) Load() *Description_ {
	s.Name().Value()
	s.Identity().Value()
	s.Website().Value()
	s.SecurityContact().Value()
	s.Details().Value()
	return s.obj
}