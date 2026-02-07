package err

import "errors"

var FailedGenerateToken = errors.New("Failed to generate Token")

var ExistingEmail = errors.New("Email already registered")

var InvalidFormatEmail = errors.New("Invalid Email Format")

var InvalidFormatPassword = errors.New("Invalid Password Format")

var ErrInvalidExt = errors.New("File have to be jpg or png")

var ErrNoRowsUpdated = errors.New("No data updated")

var WrongPassword = errors.New("Wrong password")

var WrongFormatPassword = errors.New("Wrong password format")

var MissingParameter = errors.New("Missing Parameters entered")
