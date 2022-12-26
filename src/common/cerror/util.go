package cerror

import "github.com/pkg/errors"

func CombineErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var combinedErrs error
	for _, err := range errs {
		if combinedErrs == nil && err != nil {
			combinedErrs = errors.WithStack(err)
		} else if err != nil {
			combinedErrs = errors.Wrap(combinedErrs, err.Error())
		}
	}

	return combinedErrs
}
