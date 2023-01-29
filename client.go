package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// On gère les erreurs différemment suivant si on est en connexion ou pas
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

// déclaration des constantes
const (
	IP     = "127.0.0.01" // IP local
	PORT   = "3569"       // Port utilisé
	TAILLE = 100
)

// remplissage aléatoire des fichiers décrivant les matrices
func alea(file *os.File) {
	for i := 0; i < TAILLE; i++ {
		for j := 0; j < TAILLE; j++ {
			var str = strconv.Itoa(rand.Intn(100) + 1) //converti en string un entier entre 1 et 100 inclus
			var _, err = file.WriteString(str + " ")   // écris dans le fichier
			gestionErreur(err)
		}
		var _, err = file.WriteString("\n") // écris dans le fichier
		gestionErreur(err)
	}
}

// transforme un fichier texte en chaîne de caractères
func textToString(file string) string {
	count, err := ioutil.ReadFile(file) // lire le fichier .txt
	gestionErreur(err)
	var chaine_mat = ""
	var mat = string(count)
	lines := strings.Split(mat, "\n") // on sépare les lignes du fichier
	for i := 0; i < TAILLE; i++ {
		chaine_mat += lines[i] //on sépare les caractères de chaque ligne
	}
	return chaine_mat
}

// utilisé pour lire le retour du serveur et convertir l'info en tableau d'entiers et recréer matResult
func byteToInt(b []byte, taille int) [][]int {
	var intArray [][]int             // Déclaration du tableau à double entrée d'entiers
	it := 0                          //utilisé pour parcourir le tableau d'octets
	for i := 0; i < taille; i += 1 { //parcours des lignes
		intArray = append(intArray, []int{}) //créer une nouvelle ligne au tableau
		for j := 0; j < taille; j += 1 {     //parcours des colonnes
			doubleByteVal := [][]byte{}
			var a = string(b[it : it+1])
			for a != " " { //tant que l'octet reçu ne correspond pas à un espace, stocker les octets
				doubleByteVal = append(doubleByteVal, []byte(a))
				it++
				a = string(b[it : it+1])
			}
			it++
			var byteVal = []byte{} //on regroupe tous les octets stockés dans un tableau
			for i := 0; i < len(doubleByteVal); i++ {
				byteVal = append(byteVal, doubleByteVal[i][0])
			}
			var val, err2 = strconv.Atoi(string(byteVal[:])) //transforme tout le tableau en string puis en entier, obligé de passer par string car moins complexe
			gestionErreur(err2)
			// Ajout de l'entier au sous-tableau
			intArray[i] = append(intArray[i], val) //on rajoute l'entier au bon endroit
		}
	}

	return intArray
}

// lire le buffer jusqu'au prochain espace même principe que dans byteToInt()
func lireBuff(conn net.Conn) int {
	length := [][]byte{}
	buffer_taille := make([]byte, 1)
	var _, err = conn.Read(buffer_taille)
	gestionErreurConn(err, conn)
	var b = string(buffer_taille[0])
	for b != " " {
		length = append(length, []byte(b))
		_, err = conn.Read(buffer_taille)
		b = string(buffer_taille)
	}
	var byteTaille = []byte{}
	for i := 0; i < len(length); i++ {
		byteTaille = append(byteTaille, length[i][0])
	}
	var taille, err2 = strconv.Atoi(string(byteTaille[:len(byteTaille)]))
	gestionErreur(err2)
	return taille
}

// lis une matrice dans le buffer et la converti en tableau d'entiers
func lireMat(taille int, conn net.Conn) [][]int {
	var length = lireBuff(conn) //lire la longueur en octets de la matrice dans le buffer
	buffer := make([]byte, (length))
	var _, err = io.ReadFull(conn, buffer) //lire tant qu'on n'a pas atteint la bonne longueur
	gestionErreurConn(err, conn)
	mat := byteToInt(buffer[:], taille)
	return mat
}

func main() {
	start := time.Now() //pour mesurer le temps entre le début et la fin du programme

	//Pour chaque fichier, on supprime celui existant, on en crée un nouveau vide et on le rempli d'entiers aleatoires en fonction de TAILLE
	os.Remove("A.txt")
	os.Create("A.txt")
	file_a, err := os.OpenFile("A.txt", os.O_CREATE|os.O_WRONLY, 0600)
	gestionErreur(err)
	defer file_a.Close()

	os.Remove("B.txt")
	os.Create("B.txt")
	file_b, err := os.OpenFile("B.txt", os.O_CREATE|os.O_WRONLY, 0600)
	gestionErreur(err)
	defer file_b.Close()

	os.Remove("Result.txt")
	os.Create("Result.txt")
	file_result, err := os.OpenFile("Result.txt", os.O_CREATE|os.O_WRONLY, 0600)
	gestionErreur(err)
	defer file_result.Close()

	rand.Seed(time.Now().UnixNano()) //on reset l'aléatoire pour que ça change à chque nouveau lancement de programme
	alea(file_a)
	alea(file_b)

	//On transforme les fichiers textes en chaines de caractères pour les envoyer ensuite en convertissant en tableau d'octets
	var matA = textToString("A.txt")
	var matB = textToString("B.txt")

	// Connexion au serveur
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	gestionErreur(err)

	// On envoie le message au serveur
	// le message : taille des matrices carrées, longueur en octets matrice A, matrice A, longueur en octets matrice B, matrice B
	var data = (matA)
	data = strconv.Itoa(TAILLE) + " " + strconv.Itoa(len([]byte(matA))) + " " + matA + strconv.Itoa(len([]byte(matB))) + " " + matB
	conn.Write([]byte(data))

	// On écoute le message émis par le serveur
	matResult := lireMat(TAILLE, conn)
	for i := 0; i < TAILLE; i++ {
		for j := 0; j < TAILLE; j++ {
			_, err = file_result.WriteString(strconv.Itoa(matResult[i][j]) + " ") // écrire dans le fichier
			gestionErreur(err)
		}
		var _, err = file_result.WriteString("\n") // écrire dans le fichier
		gestionErreur(err)
	}
	fmt.Println(time.Since(start))
	//une fois arrivé ici, on a fini donc on arrête la connexion avec le server
	fmt.Println("La connexion a été fermée")
	os.Exit(0)

}
