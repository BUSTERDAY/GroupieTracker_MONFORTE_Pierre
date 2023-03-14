package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type requete struct {
	name         string
	first_letter string
}

func main() {
	// Serveur pour les fichiers statiques
	fs := http.FileServer(http.Dir("HTML"))
	http.Handle("/", fs)

	css := http.FileServer(http.Dir("CSS"))
	http.Handle("/CSS/", http.StripPrefix("/CSS/", css))

	// Serveur pour les fichiers d'image
	http.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("HTML", page_requete)

	// Démarrer le serveur
	log.Fatal(http.ListenAndServe(":8080", nil))

	//récupérer les données du formulaire
	http.HandleFunc("/submit-form", submitFormHandler)
	http.ListenAndServe(":8080", nil)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)) {
	// 	if r.Method == "POST" {
	// 		name := r.FormValue("name")
	// 		ingredient := r.FormValue("ingredient")
	// 		// faire quelque chose avec les données reçues...
	// 		fmt.Fprintf(w, "Merci pour votre requête, %s ! Vous cherchez des recettes avec %s.", name, ingredient)
	// 	} else {
	// 		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	// 	}

	// 	// Construit l'URL de l'API avec le mot saisi par l'utilisateur
	// 	url_name := "www.thecocktaildb.com/api/json/v1/1/search.php?s=%s" + name
	// 	url_ingredient := "www.thecocktaildb.com/api/json/v1/1/search.php?i=%s" + ingredient
	// 	if name != "" && ingredient == "" {
	// 		url := url_name
	// 	} else if ingredient != "" && name == "" {
	// 		url := url_ingredient
	// 	}
	// }

	// // Effectue la requête GET à l'API et affiche la réponse
	// response, err := makeRequest(url)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println(response)
}

func page_requete(w http.ResponseWriter, r *http.Request) {
	// Charger la page de contact à partir d'un fichier HTML
	RequetePage, err := ioutil.ReadFile("HTML/requete.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Écrire l'en-tête de la réponse HTTP avec le type MIME correct
	w.Header().Set("Content-Type", "text/html")

	// Écrire la page de contact dans la réponse HTTP
	w.Write(RequetePage)
}

func submitFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		ingredient := r.FormValue("ingredient")
		// faire quelque chose avec les données reçues...
		fmt.Fprintf(w, "Merci pour votre requête, %s ! Vous cherchez des recettes avec %s.", name, ingredient)
	} else {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}

}
