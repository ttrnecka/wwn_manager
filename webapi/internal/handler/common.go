package handler

import "github.com/labstack/echo/v4"

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func errorWithInternal(code int, message string, err error) error {
	newErr := echo.NewHTTPError(code, message)
	newErr.Internal = err
	return newErr
}
