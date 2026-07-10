package filereader

import "errors"

var (
	ErrMalformedData   = errors.New("could not parse data to expected type")
	ErrMissingData     = errors.New("required field is empty")
	ErrInvalidFileType = errors.New("invalid filetype")
	ErrEmptyFile       = errors.New("file contain no rows")
	ErrUnknownOperator = errors.New("operator not recognized")
)
