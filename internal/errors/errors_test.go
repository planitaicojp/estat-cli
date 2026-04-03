package errors

import "testing"

func TestGetExitCode_nil(t *testing.T) {
	if got := GetExitCode(nil); got != ExitOK {
		t.Errorf("GetExitCode(nil) = %d, want %d", got, ExitOK)
	}
}

func TestGetExitCode_APIError(t *testing.T) {
	err := &APIError{StatusCode: 400, Message: "bad request"}
	if got := GetExitCode(err); got != ExitAPI {
		t.Errorf("GetExitCode(APIError) = %d, want %d", got, ExitAPI)
	}
}

func TestGetExitCode_ConfigError(t *testing.T) {
	err := &ConfigError{Message: "appId未設定"}
	if got := GetExitCode(err); got != ExitConfig {
		t.Errorf("GetExitCode(ConfigError) = %d, want %d", got, ExitConfig)
	}
}

func TestGetExitCode_NetworkError(t *testing.T) {
	err := &NetworkError{Err: nil}
	if got := GetExitCode(err); got != ExitNetwork {
		t.Errorf("GetExitCode(NetworkError) = %d, want %d", got, ExitNetwork)
	}
}

func TestGetExitCode_ValidationError(t *testing.T) {
	err := &ValidationError{Message: "invalid"}
	if got := GetExitCode(err); got != ExitValidation {
		t.Errorf("GetExitCode(ValidationError) = %d, want %d", got, ExitValidation)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 100, Code: "100", Message: "パラメータ不正"}
	got := err.Error()
	want := "APIエラー (コード 100): パラメータ不正"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
