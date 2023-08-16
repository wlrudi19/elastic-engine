package service

type Service interface {
	Sate(n, m int) int
	Bakso(n, m int) int
}

type service struct {
}

type service_2 struct {
}

func New() Service {
	return &service{}
}

func New_2() Service {
	return &service_2{}
}

func (s *service) Sate(n, m int) int {
	n, m = m, n
	return n * m
}

func (s *service) Bakso(n, m int) int {
	n, m = m, n
	return n * m
}

func (s *service_2) Sate(n, m int) int {
	n, m = m, n
	return n % m
}

func (s *service_2) Bakso(n, m int) int {
	n, m = m, n
	return n - m
}
