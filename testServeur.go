package main

import (
    "fmt"
    "net"
)

func gestionErreur(err error) {
    if err != nil {
        panic(err)
    }
}

const (
    IP   = "127.0.0.01" // IP local
    PORT = "3569"       // Port utilisé
)

func main() {

    fmt.Println("Lancement du serveur ...")

    // on écoute sur le port 3569
    ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
    gestionErreur(err)

    // On accepte les connexions entrantes sur le port 3569
    conn, err := ln.Accept()
    if err != nil {
        panic(err)
    }

    // Information sur les clients qui se connectent
    fmt.Println("Un client est connecté depuis", conn.RemoteAddr())

    gestionErreur(err)

    // boucle pour toujours écouter les connexions entrantes (ctrl-c pour quitter)
    for {
        // On écoute les messages émis par les clients
        buffer := make([]byte, 4096)       // taille maximum du message qui sera envoyé par le client
        length, err := conn.Read(buffer)   // lire le message envoyé par client
        message := string(buffer[:length]) // supprimer les bits qui servent à rien et convertir les bytes en string

        if err != nil {
            fmt.Println("Le client s'est déconnecté")
        }

        // on affiche le message du client en le convertissant de byte à string
        fmt.Print("Client:", message)

        // On envoie le message au client pour qu'il l'affiche
        conn.Write([]byte(message + "\n"))
    }
}