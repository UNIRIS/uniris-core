package contract

type ListingRepository interface {
	FindLastContract(addr string) (EndorsedContract, error)
}

type ListingService interface {
	GetLastContract(addr string) (EndorsedContract, error)
}

type listService struct {
	repo ListingRepository
}

func NewListingService(repo ListingRepository) ListingService {
	return listService{repo}
}

func (s listService) GetLastContract(addr string) (EndorsedContract, error) {
	return s.repo.FindLastContract(addr)
}
