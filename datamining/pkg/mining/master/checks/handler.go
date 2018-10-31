package checks

//Handler define methods for checkers
type Handler interface {
	CheckData(data interface{}, txHash string) error
}
