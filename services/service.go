package services

type Service struct {
	PasswdPath string
	GroupPath string
}

func NewService(passwdPath, groupPath string) *Service {

	return &Service{
		PasswdPath: passwdPath,
		GroupPath: groupPath,
	}
}