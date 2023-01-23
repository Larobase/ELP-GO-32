package main

import (
    "fmt"
    "os"
	"strings"
	"strconv"
)
func erreur(err error){
	if err != nil {
        fmt.Println(err)
        return
    }
}
