package main

import (
	"fmt"
	"net"
	"strconv"
    "io"
    "sync"
    "time"
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
	for i := 0; i < taille; i += 1 {
		intArray = append(intArray, []int{})
		for j := 0; j < taille; j += 1 {
			// Conversion des bytes en entier
			stringVal := string(b[i*taille+j : i*taille+j+1])
            intVal,err:=strconv.Atoi(stringVal)
            gestionErreur(err)
			// Ajout de l'entier au sous-tableau
			intArray[i] = append(intArray[i], intVal)
		}
	}

	return intArray
}

func multiplie(i int, a [][]int, b [][]int, taille int, conn net.Conn,c chan []byte) {
	ligne := strconv.Itoa(i)+" "
	for j := 0; j < taille; j++ {
		var mult = 0
		for k := 0; k < taille; k++ {
			mult += a[k][i] * b[j][k]
		}
        ligne += strconv.Itoa(mult)+" "
		    
	}
    ligne += "\n"
    fmt.Println(ligne)
    defer wg.Done()
	c<-[]byte(ligne)
    
}

func add(a []byte, b []byte) []byte {
	for i := 0; i < len(b); i++ {
		a = append(a, b[i])
	}
	return a
}

func client_handler(conn net.Conn) {
    // On écoute les messages émis par les clients
		buffer_taille := make([]byte, 1)
		_, err := conn.Read(buffer_taille) // lire le message envoyé par client
        taille,err :=strconv.Atoi(string(buffer_taille[:len(buffer_taille)]))
        gestionErreurConn(err,conn)
        buffer:=make([]byte,(taille*taille)*2)
        _,err=io.ReadFull(conn,buffer)
        gestionErreurConn(err,conn)
        
		matA := byteToInt(buffer[:len(buffer)/2], taille)       // supprimer les bits qui servent à rien et convertir les bytes en string
		matB := byteToInt(buffer[len(buffer)/2:len(buffer)], taille) // supprimer les bits qui servent à rien et convertir les bytes en string
        c := make(chan []byte, taille*(taille+taille))
        for i := 0; i < taille; i++ {
			wg.Add(1)
            go multiplie(i, matA, matB, taille, conn,c)
		}
        wg.Wait()
        for i := 0; i < taille; i++ {
            conn.Write(<-c)
        }
        time.Sleep(5000000000)
        // on affiche le message du client en le convertissant de byte à string
        fmt.Println("Le client s'est déconnecté")
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
