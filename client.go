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
	"bufio"
)

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	IP   = "127.0.0.01" // IP local
	PORT = "3569"       // Port utilisé
	A = 8
	B = 10
	TAILLE =3
)
func alea(file *os.File){
	for i := 0; i < TAILLE; i++ {
		for j := 0; j < TAILLE; j++ {
			var str = strconv.Itoa(rand.Intn(10) + 1)
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
	chaine_mat=strings.Replace(chaine_mat, " ", "",-1)
	return chaine_mat
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
	var result = [TAILLE]string{}

	//var matResult = [TAILLE][TAILLE]int{}
	fmt.Println(time.Since(start))
	// Connexion au serveur
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))

	gestionErreur(err)
	// On envoie le message au serveur
	time.Sleep(100)
	var data=strconv.Itoa(TAILLE)+(matA)+(matB)
	conn.Write([]byte(data))
	// On écoute tous les messages émis par le serveur et on rajouter un retour à la ligne
	for i := 0; i < TAILLE; i++ {
		var string_ligne, err = bufio.NewReader(conn).ReadString(' ')
		gestionErreur(err)
		string_ligne=string_ligne[:len(string_ligne)-1]
		valeurs, err := bufio.NewReader(conn).ReadString('\n')
		var ligne,err2=strconv.Atoi(string_ligne)
		gestionErreur(err2)
		result[ligne]=valeurs
		fmt.Print("serveur : " + valeurs)
	}
	
}