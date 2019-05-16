package main

//FeeMatrix returns the fees matrix
func FeeMatrix() (fees []float64) {
	//0.001,0.01,0.1,1,10,100,1000,10000,0.1M,1M
	for i := 0.001; i < 1000000; i *= 10 {
		fees = append(fees, i)
	}
	return
}
