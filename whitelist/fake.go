package whitelist

// FakeService ...
type FakeService struct {
	isInFunc              func() (bool, error)
	createIfNotExistsFunc func() error
	deleteFunc            func() error
}

// NewFakeService ...
func NewFakeService(isInFunc func() (bool, error), createIfNotExistsFunc func() error, deleteFunc func() error) *FakeService {
	return &FakeService{
		isInFunc:              isInFunc,
		createIfNotExistsFunc: createIfNotExistsFunc,
		deleteFunc:            deleteFunc,
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

// Delete ...
func (s *FakeService) Delete(id int64) error {
	return s.deleteFunc()
}
