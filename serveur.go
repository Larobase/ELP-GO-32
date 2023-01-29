package main

import (
	"fmt"
	"net"
	"strconv"
    "io"
    "sync"
)
// On gère les erreurs différemment suivant si on est en connexion ou pas
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

var wg sync.WaitGroup //on déclare ici  pour y avoir accès dans les fonctions

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

//fonction exécutée simultanément pour chaque ligne du tableau avec des go routines
func multiplie(i int, a , b, r *[][]int, taille int) {
	for j := 0; j < taille; j++ {
		var mult = 0
		for k := 0; k < taille; k++ {
			mult += (*a)[i][k] * (*b)[k][j] 
		}
		(*r)[i][j]=mult	    
	}
	defer wg.Done()	//on indique au waitgroup que le travil est terminé
}

//transforme un tableau d'entiers en une chaine de caractères
func intToString(r [][]int) string{
	var str = ""
	for i := 0; i < len(r); i++ {
		for j := 0; j < len(r[i]); j++ {
			str+=strconv.Itoa(r[i][j])+" " //conversion de int vers string
		}
	}
	return str
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

// gère le client qui vient de se connecter
func client_handler(conn net.Conn) {
    // On écoute les messages émis par les clients
	var taille=lireBuff(conn)	//on récupère taille des matrices carrées
	matA := lireMat(taille,conn)
	matB := lireMat(taille,conn)	
	matResult := make([][]int, taille) //on créé la matrice résultat
	for i:=0;i<len(matResult);i++ {
		matResult[i] = make([]int, taille) //on lui ajoute autant de cases que matA et matB
	}
    for i := 0; i < taille; i++ {
		wg.Add(1)		//on rajoute un worker
        go multiplie(i, &matA, &matB,&matResult, taille) //on fait le calcul des lignes en simultané
	}
    wg.Wait() //on attend que toutes les lignes aient fini de se calculer
	var data=intToString(matResult) 
	data = strconv.Itoa(len([]byte(data)))+" "+data
	conn.Write([]byte(data)) //on renvoie la matrice résultat
    fmt.Println("Le client s'est déconnecté") //on déconnecte le client une fois que l'envoie est fait
	conn.Close()
              
		
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
