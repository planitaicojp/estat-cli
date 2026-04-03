package errors

import "fmt"

type ExitCoder interface {
	ExitCode() int
}

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("設定エラー: %s", e.Message)
}

func (e *ConfigError) ExitCode() int {
	return ExitConfig
}

type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIエラー (コード %s): %s", e.Code, e.Message)
}

func (e *APIError) ExitCode() int {
	return ExitAPI
}

type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("ネットワークエラー: %v", e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

func (e *NetworkError) ExitCode() int {
	return ExitNetwork
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("バリデーションエラー (%s): %s", e.Field, e.Message)
	}
	return fmt.Sprintf("バリデーションエラー: %s", e.Message)
}

func (e *ValidationError) ExitCode() int {
	return ExitValidation
}

func GetExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitGeneral
}
