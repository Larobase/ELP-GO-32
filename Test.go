package main

import (
	"fmt"
	"math/rand"

)

const TAILLEMAT=2

func multiplie(i int, j int,a *[][]int,b*[][]int,c chan int) int{
	mult := 0
	for k := 0; k < len(*a); k++ {
		mult += (*a)[k][i]*(*b)[j][k]
	}
	return mult
	//c <- mult // send mult to c
}
func max(a int, b int) int {
	if (a>b){
		return a
	}
	return b

	
}

func colonne(a[][]int,index int,) []int{
	column := make([]int, len(a))
    for i := range a {
        column[i] = a[i][index]
    }
    return column

}

func main() {
	a:= [][]int{}
	b:= [][]int{}
	for i := 0; i < TAILLEMAT; i++ {
		rowa := make([]int, TAILLEMAT)
		rowb := make([]int, TAILLEMAT)
		for j := 0; j < TAILLEMAT; j++ {
			rowa[j] =rand.Intn(10) + 1
			rowb[j] =rand.Intn(10) + 1
		}
		a= append(a, rowa)
		b= append(b, rowb)

	}
	result := [][]int{}
	for i :=  range a {
		row := make([]int, len(a))
		for j :=  range b[0] {
			row[j] = (i+1)/(i+1)-1
		}
		result = append(result, row)
	}
	
	fmt.Println(a)
	fmt.Println(b)
	
	for i := 0; i < len(a); i++ {
			for j := 0; j < len(a[0]); j++ {
				result[i][j]= multiplie(i,j,&a,&b)
		}
	}
	fmt.Println(result)
	

} 

