package cardano

import "fmt"

// error when failing to parse time entities.
type parsingError struct {
    // text that should have been parsed.
    ParsedText string
    // reason why it cannot be parsed.
    Reason string
}

func (err parsingError) Error() string {
    return fmt.Sprintf("Failed to parse '%v'. %v", err.ParsedText, err.Reason)
}

type invalidArgument struct {
    MethodName string
    Expected   string
}

func (err invalidArgument) Error() string {
    return fmt.Sprintf("Invalid argument passed to '%v'. %v", err.MethodName, err.Expected)
}
