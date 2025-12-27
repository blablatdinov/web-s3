package srv_test

import (
	"testing"

	repo "github.com/blablatdinov/web-s3/src/repo"
	srv "github.com/blablatdinov/web-s3/src/srv"
)

func Test(t *testing.T) {
	pswrdHash, err := srv.PswrdCtor("fkPassword").Hash()
	if err != nil {
		t.Fatalf("Fail to generate hash: %s", err.Error())
	}
	authSrv := srv.UserAuthSrvCtor(
		"fkSecret",
		repo.FkUserAuthRepoCtor(
			0,
			pswrdHash,
		),
	)
	token, err := authSrv.Jwt("user1", "fkPassword")
	if err != nil {
		t.Fatalf("Fail on generate token: %s", err.Error())
	}
	valid, err := authSrv.Validate(token)
	if err != nil {
		t.Fatalf("Fail on validate token: %s", err.Error())
	}
	if !valid {
		t.Fatalf("Invalid token")
	}
}
