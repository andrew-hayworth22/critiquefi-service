package postgres

type sysPG struct {
	repo *repo
}

func (s *sysPG) Ping() error {
	return s.repo.Ping()
}
