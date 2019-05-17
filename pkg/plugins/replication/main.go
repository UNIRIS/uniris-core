package main

func GetMinimumReplicas(tx interface{}) int {
	return 1
}

func IsAuthorizedToStoreTx(tx interface{}) bool {
	return true
}
