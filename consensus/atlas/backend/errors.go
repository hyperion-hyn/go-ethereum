package backend

import "errors"

var (
    // errFailedSignData is returned when failed to sign data
    errFailedSignData = errors.New("failed to sign data")
    // errValidatorNotExist is return if no validator in statedb.
    errValidatorNotExist = errors.New("staking validator does not exist")
)
