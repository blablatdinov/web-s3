package srv_test

import (
	"errors"
	"testing"

	"github.com/blablatdinov/web-s3/src/repo"
	srv "github.com/blablatdinov/web-s3/src/srv"
)

func TestEmptyUsername(t *testing.T) {
	usrSignupSrv := srv.UsrSignupSrvCtor(repo.FkUserSignupRepoCtor(0, nil))
	_, err := usrSignupSrv.Create("", "pass")
	if err == nil {
		t.Fatalf("Error not raised")
	}
	if !errors.Is(err, srv.ErrUsernameEmpty) {
		t.Fatalf("Error not matched")
	}
}
