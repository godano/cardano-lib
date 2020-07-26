package time

import "fmt"

// error when failing to parse time entities.
type ParsingError struct {
    // text that should have been parsed.
    ParsedText string
    // reason why it cannot be parsed.
    Reason string
}

func (err ParsingError) Error() string {
    return fmt.Sprintf("Failed to parse '%v'. %v", err.ParsedText, err.Reason)
}

type InvalidArgument struct {
    MethodName string
    Expected   string
}

func (err InvalidArgument) Error() string {
    return fmt.Sprintf("Invalid argument passed to '%v'. %v", err.MethodName, err.Expected)
}
