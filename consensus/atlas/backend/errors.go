package backend

import "errors"

var (
	// errFailedSignData is returned when failed to sign data
	errFailedSignData = errors.New("failed to sign data")
	// errValidatorNotExist is return if no validator in statedb.
	errValidatorNotExist       = errors.New("validator does not exist")
	errRedelegationNotExist    = errors.New("redelegation does not exists")
	errMap3NodeNotExist        = errors.New("map3 node does not exist")
	errMicrodelegationNotExist = errors.New("microdelegation does not exists")
)
