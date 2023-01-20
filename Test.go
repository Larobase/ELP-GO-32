package main

import "fmt"

func multiplie(i int,j int, a[][]int,b[][]int,r[][]int,c chan int,) {
	mult := 0
	for k := 0; k < len(a); k++ {
		r[i][j] += a[i][k]*b[k][j]
	}
	c <- mult // send mult to c
}

func main() {
	a := [][]int{{2,2,2},{2,2,2}}	
	b := [][]int{{2,2,2},{2,2,2}}
	 r [2][3]int
	c := make(chan int)
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a); j++ {
			go multiplie(i,j,a,b,r,c)
		}
	}
	fmt.Println(r)

} 

