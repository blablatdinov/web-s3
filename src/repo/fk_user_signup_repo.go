package repo

type FkUserSignupRepo struct {
	userId int
	err    error
}

func FkUserSignupRepoCtor(userId int, err error) UserSignupRepo {
	return FkUserSignupRepo{}
}

func (r FkUserSignupRepo) Create(username, passwordHash string) (int, error) {
	return r.userId, r.err
}
