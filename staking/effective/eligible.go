package effective

import "github.com/ethereum/go-ethereum/staking/types/restaking"

// Candidacy is a more semantically meaningful
// value that is derived from core protocol logic but
// meant more for the presentation of user, like at RPC
type Candidacy byte

const (
	// Unknown ..
	Unknown Candidacy = iota
	// ForeverBanned ..
	ForeverBanned
	// Candidate ..
	Candidate = iota
	// NotCandidate ..
	NotCandidate
	// Elected ..
	Elected
)

const (
	doubleSigningBanned = "banned forever from network because was caught double-signing"
)

func (c Candidacy) String() string {
	switch c {
	case ForeverBanned:
		return doubleSigningBanned
	case Candidate:
		return "eligible to be elected next epoch"
	case NotCandidate:
		return "not eligible to be elected next epoch"
	case Elected:
		return "currently elected"
	default:
		return "unknown"
	}
}

// ValidatorStatus ..
func ValidatorStatus(currentlyInCommittee bool, status restaking.ValidatorStatus) Candidacy {
	switch {
	case status == restaking.Banned:
		return ForeverBanned
	case currentlyInCommittee:
		return Elected
	case !currentlyInCommittee && status == restaking.Active:
		return Candidate
	case !currentlyInCommittee && status != restaking.Active:
		return NotCandidate
	default:
		return Unknown
	}
}

// BootedStatus ..
type BootedStatus byte

const (
	// Booted ..
	Booted BootedStatus = iota
	// NotBooted ..
	NotBooted
	// LostEPoSAuction ..
	LostEPoSAuction
	// TurnedInactiveOrInsufficientUptime ..
	TurnedInactiveOrInsufficientUptime
	// BannedForDoubleSigning ..
	BannedForDoubleSigning
)

func (r BootedStatus) String() string {
	switch r {
	case Booted:
		return "booted"
	case LostEPoSAuction:
		return "lost epos auction"
	case TurnedInactiveOrInsufficientUptime:
		return "manually turned inactive or insufficient uptime"
	case BannedForDoubleSigning:
		return doubleSigningBanned
	default:
		return "not booted"
	}
}
