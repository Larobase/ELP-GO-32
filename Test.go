package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
	fmt.Print("hola")
}

const (
	IP   = "127.0.0.01" // IP local
	PORT = "3569"       // Port utilisé
)

func main() {

	// Connexion au serveur
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	gestionErreur(err)

	for {
		// entrée utilisateur
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("client: ")
		text, err := reader.ReadString('\n')
		gestionErreur(err)

		// On envoie le message au serveur
		conn.Write([]byte(text))

		// On écoute tous les messages émis par le serveur et on rajouter un retour à la ligne
		message, err := bufio.NewReader(conn).ReadString('\n')
		gestionErreur(err)

		// on affiche le message utilisateur
		fmt.Print("serveur : " + message)
	}
}