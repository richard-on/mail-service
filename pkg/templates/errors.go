package templates

import "errors"

var ErrBadFormat = errors.New("unable to process: incorrect or missing keys")

var ErrNoSuchTemplate = errors.New("unknown template type")
