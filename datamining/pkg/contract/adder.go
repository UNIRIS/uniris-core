package contract

type AddingRepository interface {
	StoreEndorsedContract(EndorsedContract) error
}

type AddingService interface {
	StoreEndorsedContract(EndorsedContract) error
}

type addService struct {
	repo AddingRepository
}

func NewAddingService(repo AddingRepository) AddingService {
	return addService{repo}
}

func (s addService) StoreEndorsedContract(c EndorsedContract) error {
	//REMOVE STORAGE CHECKS FOR DEMO PURPOSE

	return s.repo.StoreEndorsedContract(c)
}
