package osb

import "github.com/pkg/errors"

var (
	ErrInstanceIDNotFound = errors.New("the specified instance_id does not exist")
)
