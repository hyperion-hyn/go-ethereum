package restaking

import "github.com/pkg/errors"

var (
	MaxNameLength            = 140
	MaxIdentityLength        = 140
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

func (d *Description_) IncrementalUpdateFrom(other Description_) error {
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
	if len(d.Identity) == 0 {
		return errors.New("empty identity")
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

func (s *Storage_Description_) IncrementalUpdateFrom(other Description_) {
	if other.Name != "" {
		s.Name().SetValue(other.Name)
	}
	if other.Identity != "" {
		s.Identity().SetValue(other.Identity)
	}
	if other.Website != "" {
		s.Website().SetValue(other.Website)
	}
	if other.SecurityContact != "" {
		s.SecurityContact().SetValue(other.SecurityContact)
	}
	if other.Details != "" {
		s.Details().SetValue(other.Details)
	}
}
