package main

import (
	"fmt"
	"net"
	"strconv"
    "io"
    "sync"
)

func gestionErreurConn(err error, conn net.Conn) {
	if err != nil {
        fmt.Println("Le client s'est déconnecté")
		conn.Close()
	}
}
func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	IP   = "127.0.0.01" // IP local
	PORT = "3569"       // Port utilisé
)

var wg sync.WaitGroup

func byteToInt(b []byte, taille int) [][]int {
	// Déclaration du tableau à double entrée d'entiers
	var intArray [][]int
	// Boucle sur chaque 4 bytes pour décoder les entiers
	it:=0
	for i := 0; i < taille; i += 1 {
		intArray = append(intArray, []int{})
		for j := 0; j < taille; j += 1 {
			doubleByteVal:=[][]byte{}
			var a= string(b[it: it+1])
			for a !=" "{
				doubleByteVal=append(doubleByteVal,[]byte(a))
				it++
				a=string(b[it: it+1])
			}
			it++
			var byteVal = []byte{}
			for i:=0;i<len(doubleByteVal);i++{
				byteVal=append(byteVal, doubleByteVal[i][0])
			}
			var val,err2=strconv.Atoi(string(byteVal[:]))
			gestionErreur(err2)
			// Ajout de l'entier au sous-tableau
			intArray[i] = append(intArray[i], val)
		}
	}

	return intArray
}

func multiplie(i int, a , b, r *[][]int, taille int) {
	for j := 0; j < taille; j++ {
		var mult = 0
		for k := 0; k < taille; k++ {
			mult += (*a)[i][k] * (*b)[k][j]
		}
		(*r)[i][j]=mult	    
	}
	defer wg.Done()
}

func add(a []byte, b []byte) []byte {
	for i := 0; i < len(b); i++ {
		a = append(a, b[i])
	}
	return a
}

func intToString(r [][]int) string{
	var str = ""
	for i := 0; i < len(r); i++ {
		for j := 0; j < len(r[i]); j++ {
			str+=strconv.Itoa(r[i][j])+" "
		}
	}
	return str
}

func lireBuff(conn net.Conn)int{
	length:=[][]byte{}
		buffer_taille := make([]byte, 1)
		var _,err=conn.Read(buffer_taille)
		gestionErreurConn(err,conn)
		var b=string(buffer_taille[0])
		for b !=" "{
			length=append(length,[]byte(b))
			_,err=conn.Read(buffer_taille)
			 b=string(buffer_taille)
		}
		var byteTaille = []byte{}
		for i:=0;i<len(length);i++{
			byteTaille=append(byteTaille, length[i][0])
		}
		var taille,err2=strconv.Atoi(string(byteTaille[:len(byteTaille)]))
		gestionErreur(err2)
		return taille	
}

func lireMat(taille int,conn net.Conn)[][]int{
	var length=lireBuff(conn)
	buffer:=make([]byte,(length))
    var _,err=io.ReadFull(conn,buffer)
    gestionErreurConn(err,conn)
	mat := byteToInt(buffer[:], taille)
	return mat
}

func client_handler(conn net.Conn) {
    // On écoute les messages émis par les clients
	var taille=lireBuff(conn)	
	matA := lireMat(taille,conn)
	matB := lireMat(taille,conn)	
	matResult := make([][]int, taille)
	for i:=0;i<len(matResult);i++ {
		matResult[i] = make([]int, taille)
	}
    for i := 0; i < taille; i++ {
		wg.Add(1)
        go multiplie(i, &matA, &matB,&matResult, taille)
	}
    wg.Wait()
	var data=intToString(matResult)
	data = strconv.Itoa(len([]byte(data)))+" "+data
	conn.Write([]byte(data))
    // on affiche le message du client en le convertissant de byte à string
    fmt.Println("Le client s'est déconnecté")
              
		
}

func main() {

	fmt.Println("Lancement du serveur ...")

	// on écoute sur le port 3569
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	gestionErreur(err)

	// boucle pour toujours écouter les connexions entrantes (ctrl-c pour quitter)
	for {
		// On accepte les connexions entrantes sur le port 3569
		conn, err := ln.Accept()
		gestionErreurConn(err, conn)

		// Information sur les clients qui se connectent
		fmt.Println("Un client est connecté depuis", conn.RemoteAddr())
		gestionErreurConn(err, conn)
		go client_handler(conn)
		
	}
}
