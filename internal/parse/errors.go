package parse

import (
	"fmt"
	"log"
)

func CombineErrs(args ...any) error {
	// args: label1, []error1, label2, []error2, ...
	var outErr error
	for i := 0; i < len(args); i += 2 {
		label := args[i].(string)
		errs := args[i+1].([]error)
		if len(errs) == 0 {
			continue
		}
		// build a single error with a short summary
		outErr = fmt.Errorf("%s: %d parse errors (first: %v)", label, len(errs), errs[0])
	}
	return outErr
}

func ErrorPrinter(errs []error, limit int) {
	if len(errs) < limit {
		limit = len(errs)
	}

	for _, err := range errs[:limit] {
		log.Printf("parsing errors: %s", err)
	}
}
