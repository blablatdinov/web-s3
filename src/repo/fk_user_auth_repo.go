package repo

type FkUserAuthRepo struct {
	userId       int
	passwordHash string
}

func FkUserAuthRepoCtor(UserId int, PasswordHash string) UserAuthRepo {
	return FkUserAuthRepo{UserId, PasswordHash}
}

func (repo FkUserAuthRepo) UserId(username string) (int, error) {
	return repo.userId, nil
}

func (repo FkUserAuthRepo) PasswordHash(username string) (string, error) {
	return repo.passwordHash, nil
}
