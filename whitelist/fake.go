package whitelist

// FakeService ...
type FakeService struct {
	isInFunc              func() (bool, error)
	createIfNotExistsFunc func() error
}

// NewFakeService ...
func NewFakeService(isInFunc func() (bool, error), createIfNotExistsFunc func() error) *FakeService {
	return &FakeService{
		isInFunc:              isInFunc,
		createIfNotExistsFunc: createIfNotExistsFunc,
	}
}

// IsIn ...
func (s *FakeService) IsIn(id int64) (bool, error) {
	return s.isInFunc()
}

// CreateIfNotExists ...
func (s *FakeService) CreateIfNotExists(id int64, username string) error {
	return s.createIfNotExistsFunc()
}
