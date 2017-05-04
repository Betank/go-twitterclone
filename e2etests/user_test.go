package e2etests

import "testing"

func TestCreateUser(t *testing.T) {
	err := registerNewUser("newUser", "test")
	if err != nil {
		t.Error(err)
	}
}

func TestCreateUserAlreadyExists(t *testing.T) {
	err := registerNewUser("newUser2", "test")
	if err != nil {
		t.Error(err)
	}

	err = registerNewUser("newUser2", "test")
	if err != userConflictError {
		t.Error("should fail with status conflict")
	}
}

func TestLogin(t *testing.T) {
	err := registerNewUser("user", "test")
	if err != nil {
		t.Error(err)
	}

	token, err := userLogin("user", "test")
	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error("token is empty")
	}
}

func TestLoginFailed(t *testing.T) {
	err := registerNewUser("user2", "test")
	if err != nil {
		t.Error(err)
	}

	token, err := userLogin("user2", "test2")
	if err != loginFailedError {
		t.Error("should fail with status unauthorized")
	}

	if token != "" {
		t.Error("token should be empty")
	}
}
