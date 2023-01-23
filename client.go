package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
	"math/rand"
    "io/ioutil"	
	"time"
	"encoding/binary"
)

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
	fmt.Print("ceci est un test d'erreur\n")
}
func gestionErreur2(err error) {
	if err != nil {
		panic(err)
	}
	fmt.Print("erreur recep\n")
}

const (
	IP   = "127.0.0.01" // IP local
	PORT = "3569"       // Port utilisé
	A = 8
	B = 10
	TAILLE =1
)
func alea (file *os.File){
	for i := 0; i < TAILLE; i++ {
		for j := 0; j < TAILLE; j++ {
			var str = strconv.Itoa(rand.Intn(10) + 1)
			var _, err = file.WriteString(str) // écrire dans le fichier
			gestionErreur(err)
		}
		if i!=(TAILLE-1) {
			var _, err = file.WriteString("\n") // écrire dans le fichier
			gestionErreur(err)
		}
	}
}

func intToByte(a [][]int) []byte{
	// Déclaration d'un buffer pour stocker les données encodées
	var buf []byte

	// Boucle sur chaque sous-tableau pour encoder les entiers
	for _, subArray := range a {
		for _, intVal := range subArray {
			// Encodage de l'entier en bytes
			var intBytes [4]byte
			binary.LittleEndian.PutUint32(intBytes[:], uint32(intVal))
			// Ajout des bytes encodés au buffer
			buf = append(buf, intBytes[:]...)
		}
	}
	return buf
}

func add(a []byte,b []byte)[]byte{
	for i := 0; i < len(b); i++ {
		a=append(a,b[i])
	}
	return a
}

func extraction(file string) [][]int{
    count, err := ioutil.ReadFile(file) // lire le fichier text.txt
    gestionErreur(err)

    var mat = string(count)
	lines := strings.Split(mat, "\n")
	chars := make([][]string, len(lines))
	number := make([][]int, len(lines))
	var nb=0
	for i := 0; i < TAILLE; i++ {
		chars[i] = strings.Fields(lines[i])
	}
	for i := 0; i < TAILLE; i++ {
		for j := 0; j < TAILLE; j++ {
			nb,err = strconv.Atoi(chars[i][j])
			gestionErreur2(err)
			number[i] = append(number[i],nb)

		}
	}
	return number
}

func main() {
	start :=time.Now()
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
	gestionErreur2(err)
	defer file_result.Close()

	rand.Seed(time.Now().UnixNano())
	alea(file_a)
	alea(file_b)

	var matA = extraction("A.txt")
	var matB = extraction("B.txt")
	fmt.Println(time.Since(start))
	
	// Connexion au serveur
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	gestionErreur(err)

	for {
		
		// On envoie le message au serveur
		time.Sleep(100)
		fmt.Print("client: ",matA,matB)
		var data=intToByte(matA)
		data=add(data,intToByte(matB))
		conn.Write(data)

		// On écoute tous les messages émis par le serveur et on rajouter un retour à la ligne
		message,err:= bufio.NewReader(conn).ReadBytes(100)
		gestionErreur(err)
		fmt.Print("ceci est un test\n")
		var envoie=[3]int{}
		envoie[0] = int(binary.LittleEndian.Uint32(message[0:4]))
		envoie[1] = int(binary.LittleEndian.Uint32(message[4:8]))
		envoie[2] = int(binary.LittleEndian.Uint32(message[8:len(message)]))
		fmt.Print("serveur : " + string(envoie[2]))
	}
}