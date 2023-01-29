package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
	"math/rand"
    "io/ioutil"	
	"time"
	"io"
)

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}

func gestionErreurConn(err error, conn net.Conn) {
	if err != nil {
        fmt.Println("Le serveur s'est déconnecté")
		conn.Close()
	}
}


const (
	IP   = "127.0.0.01" // IP local
	PORT = "3569"       // Port utilisé
	TAILLE =100
)
func alea(file *os.File){
	for i := 0; i < TAILLE; i++ {
		for j := 0; j < TAILLE; j++ {
			var str = strconv.Itoa(rand.Intn(100) + 1)
			var _, err = file.WriteString(str+" ") // écrire dans le fichier
			gestionErreur(err)
		}
		var _, err = file.WriteString("\n") // écrire dans le fichier
		gestionErreur(err)
	}
}

func textToString(file string) string{
    count, err := ioutil.ReadFile(file) // lire le fichier text.txt
    gestionErreur(err)
	var chaine_mat=""
    var mat = string(count)
	lines := strings.Split(mat, "\n")
	for i := 0; i < TAILLE; i++ {
		chaine_mat+=lines[i]
	}
	return chaine_mat
}

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

func main() {
	start :=time.Now() //pour mesurer le temps entre le début et la fin du programme
//Pour chaque fichier, on supprime celui existant, on en crée un nouveau vide et on le rempli d'entiers aleatoires en fonction de TAILLE
	os.Remove("A.txt")
	os.Create("A.txt")
	file_a,err := os.OpenFile("A.txt", os.O_CREATE|os.O_WRONLY, 0600)
	gestionErreur(err)
	defer file_a.Close()

	os.Remove("B.txt")
	os.Create("B.txt")
	file_b,err := os.OpenFile("B.txt", os.O_CREATE|os.O_WRONLY, 0600)
	gestionErreur(err)
	defer file_b.Close()

	os.Remove("Result.txt")
	os.Create("Result.txt")
	file_result,err := os.OpenFile("Result.txt", os.O_CREATE|os.O_WRONLY, 0600)
	gestionErreur(err)
	defer file_result.Close()

	rand.Seed(time.Now().UnixNano()) //on reset l'aléatoire pour que ça change à chque nouveau lancement de programme
	alea(file_a)
	alea(file_b)

	//On transforme les fichiers textes en 
	var matA = textToString("A.txt")
	var matB = textToString("B.txt")
	
	// Connexion au serveur
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	gestionErreur(err)
	// On envoie le message au serveur
	var data=(matA)
	data=strconv.Itoa(TAILLE)+" "+strconv.Itoa(len([]byte(matA)))+" "+matA+strconv.Itoa(len([]byte(matB)))+" "+matB
	conn.Write([]byte(data))
	// On écoute tous les messages émis par le serveur et on rajouter un retour à la ligne
	matResult := lireMat(TAILLE,conn)
	for i := 0; i < TAILLE; i++ {
		for j := 0; j < TAILLE; j++ {
			var _, err = file_result.WriteString(strconv.Itoa(matResult[i][j])+" ") // écrire dans le fichier
			gestionErreur(err)
		}
		var _, err = file_result.WriteString("\n") // écrire dans le fichier
		gestionErreur(err)
	}
	fmt.Println(time.Since(start))
	fmt.Println("Le serveur a fermé la connexion")
	os.Exit(0)
	
}