package main

import (
    "fmt"
    "net"
    "encoding/binary"
    "math"
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

func byteToInt(b []byte,taille int) [][]int {
    // Déclaration du tableau à double entrée d'entiers
	var intArray [][]int
    // Boucle sur chaque 4 bytes pour décoder les entiers
	for i := 0; i < taille ; i += 1 {
        intArray = append(intArray, []int{})
        for j := 0; j < taille ; j += 1 {
            // Conversion des bytes en entier
		    intVal := int(binary.LittleEndian.Uint32(b[4*i*taille+4*j : 4*i*taille+4*j+4]))
		    // Ajout de l'entier au sous-tableau
		    intArray[i] = append(intArray[i], intVal)
        }
	}
    return intArray
}

func multiplie(i int,j int,a*[]int, b *[]int,taille int) int {
    var mult = 0
	for k := 0; k < taille; k++ {
		mult+= (*a)[i]*(*b)[j]
	}
	return mult
}

func add(a []byte,b []byte)[]byte{
	for i := 0; i < len(b); i++ {
		a=append(a,b[i])
	}
	return a
}

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
	var end=true
    // boucle pour toujours écouter les connexions entrantes (ctrl-c pour quitter)
    for {
        // On écoute les messages émis par les clients
        buffer := make([]byte,4096)
        length, err := conn.Read(buffer)   // lire le message envoyé par client
        message := (buffer[:length])
        var taille = int(math.Sqrt(float64(len(message)/8)))
        matA := byteToInt(buffer[:length/2],taille) // supprimer les bits qui servent à rien et convertir les bytes en string
        matB := byteToInt(buffer[length/2:length],taille) // supprimer les bits qui servent à rien et convertir les bytes en string
        
        if message!=[]int{
			end = true
		}
        if err != nil{
            fmt.Println("Le client s'est déconnecté")
			end=false
        }
		
		if end==true{
			for i := 0; i < taille; i++ {
                for j := 0; j < taille; j++ {
                    var calcul = multiplie(i,j,&matA[i],&matB[:len(matB[0])][j],taille)
                    envoie:=[]int{i,j,calcul}
                     intBytes := make([]byte,4096)
                     passage :=make([]byte,4096)
                    for k := 0; k < len(envoie); k++ {
                        binary.LittleEndian.PutUint32(passage, uint32(envoie[k]))
                        intBytes=add(intBytes,passage)
                    }
                    conn.Write(intBytes)
                }
  
            }
            // on affiche le message du client en le convertissant de byte à string
			 fmt.Println("Client:", matA)
             fmt.Println("Client:", matB)
		}
    }
}