package filereader

import "errors"

var (
	ErrMalformedData   = errors.New("could not convert parse data to expected type")
	ErrMissingData     = errors.New("required field is empty")
	ErrInvalidFileType = errors.New("invalid filetype")
	ErrEmptyFile       = errors.New("file contain no rows")
)
