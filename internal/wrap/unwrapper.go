package wrap

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func UnwrapDiscordRESTErr(err error, allowsCodes ...int) error {
	if err == nil {
		return err
	}
	mappedCodes := make(map[int]bool)
	for _, code := range allowsCodes {
		mappedCodes[code] = true
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("UnwrapErr", err))
		}
	}()

	var discordGoErr *discordgo.RESTError
	if errors.As(err, &discordGoErr) {
		if mappedCodes[discordGoErr.Message.Code] {
			return nil
		}
	}
	return err
}

func UnwrapErrToErrSlices(err error) []error {
	var errs []error
	if unwrappedErrors, ok := err.(interface{ Unwrap() []error }); ok {
		for _, unwrappedErr := range unwrappedErrors.Unwrap() {
			errs = append(errs, UnwrapErrToErrSlices(unwrappedErr)...)
		}
	} else {
		errs = append(errs, err)
	}
	return errs
}
