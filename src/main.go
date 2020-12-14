package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	port := "8080" /*Changer le port du serveur ici*/
	fmt.Println("Starting server on 127.0.0.1:" + port + "/")
	http.HandleFunc("/", index)
	http.HandleFunc("/download", downloadAscii)
	http.ListenAndServe(":"+port, nil)
}

/*Route Principale*/
func index(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	t = template.Must(template.ParseGlob("templates/index.html"))
	if r.URL.Path != "/" {
		http.ServeFile(w, r, "templates/404.html")
		return
	}
	var text, font, Result string
	switch r.Method {
	case "GET":
		t.ExecuteTemplate(w, "index", nil)
	case "POST":
		font = r.FormValue("form")
		text = r.FormValue("textInput")

		argRune := []rune(text)
		for index1 := 0; index1 < len(argRune); index1++ {
			if argRune[index1] < 32 || argRune[index1] > 126 {
				http.ServeFile(w, r, "templates/400.html")
				return
			}
		}
		print(text, font)
	}
	content, err := ioutil.ReadFile("file.txt")
	if err != nil {
		http.ServeFile(w, r, "templates/500.html")
	}
	tabContent := []byte(content)
	Result = string(tabContent[1:])
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, Result)
}

/*Route pour télécharger*/
func downloadAscii(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "file.txt")
}

//fonction pour print
func print(Args string, Style string) {
	Arg := []rune(Args)
	var nbligne, ligne, charRetour int
	debutMot := 0
	retour := false
	var forms string
	if Style == "standard" {
		forms = "forms/standard.txt"
	} else if Style == "shadow" {
		forms = "forms/shadow.txt"
	} else {
		forms = "forms/thinkertoy.txt"
	}
	file, err := os.Open("file.txt")
	if err == nil {
		file.Close()
		os.Remove("file.txt")
	}
	nbretour := retourLigne(Arg)
	for retourLigne := 0; retourLigne <= nbretour; retourLigne++ {
		for lignePrint := 0; lignePrint < 8; lignePrint++ {
			for lettre := debutMot; lettre < len(Arg); lettre++ {
				nbligne = 0
				ligne = getLigne(Arg[lettre])
				file, erreur := os.Open(forms)
				if erreur != nil {
					error(erreur.Error())
				} else {
					scanner := bufio.NewScanner(file)
					if Arg[lettre] == '\\' && Arg[lettre+1] == 'n' {
						retour = true
						charRetour = lettre + 2
						break
					}
					for scanner.Scan() {
						if nbligne == ligne+lignePrint {
							file, err := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
							defer file.Close()
							if err != nil {
								error(err.Error())
							}
							_, err = file.WriteString(scanner.Text()) // écrire dans le fichier
							if err != nil {
								error(err.Error())
							}
						}
						nbligne++
					}
				}
			}
			file, _ := os.OpenFile("file.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			_, err := file.WriteString("\n")
			if err != nil {
				error(err.Error())
			}
			file.Close()
		}
		if retour {
			debutMot = charRetour
		}
	}
}

//Fonction récupérer la ligne du symbole
func getLigne(char rune) int {
	var ligne int
	for index := 0; index < 95; index++ {
		if rune(index+32) == char {
			ligne = index
			break
		}
	}
	ligne = ligne*9 + 1
	return ligne
}

//Fonction pour recuperer le nombre de retour à la ligne
func retourLigne(Arg []rune) (nbretour int) {
	count := 0
	for index := 0; index < len(Arg); index++ {
		if Arg[index] == '\\' && Arg[index+1] == 'n' {
			count++
		}
	}
	return count
}

//Fonction erreur qui affiche l'erreur
func error(str string) {
	fmt.Println("ERREUR: " + str)
}
