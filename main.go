package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

	// Écrire la page de requete dans la réponse HTTP
	w.Write(RequetePage)
}

type Cocktail struct {
	Name           string `json:"strDrink"`
	Instructions   string `json:"strInstructions"`
	ThumbnailImage string `json:"strDrinkThumb"`
}

func submitFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		ingredient := r.FormValue("ingredient")

		// Créez l'URL de recherche de cocktail en utilisant les paramètres d'entrée
		searchByNameURL := fmt.Sprintf("https://www.thecocktaildb.com/api/json/v1/1/search.php?s=%s", url.QueryEscape(name))
		searchByIngredientURL := fmt.Sprintf("https://www.thecocktaildb.com/api/json/v1/1/filter.php?i=%s", url.QueryEscape(ingredient))

		// Effectuez une requête GET pour les cocktails correspondants
		searchByNameRes, err := http.Get(searchByNameURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer searchByNameRes.Body.Close()

		searchByIngredientRes, err := http.Get(searchByIngredientURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer searchByIngredientRes.Body.Close()

		// Parsez les résultats de la recherche en JSON
		var searchByNameData map[string]interface{}
		err = json.NewDecoder(searchByNameRes.Body).Decode(&searchByNameData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var searchByIngredientData map[string]interface{}
		err = json.NewDecoder(searchByIngredientRes.Body).Decode(&searchByIngredientData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Récupérez les résultats de la recherche
		var cocktails []Cocktail

		// Si la recherche par nom de cocktail est fructueuse, récupérez les cocktails correspondants
		if searchByNameData["drinks"] != nil {
			for _, drink := range searchByNameData["drinks"].([]interface{}) {
				d := drink.(map[string]interface{})
				cocktails = append(cocktails, Cocktail{
					Name:           d["strDrink"].(string),
					Instructions:   d["strInstructions"].(string),
					ThumbnailImage: d["strDrinkThumb"].(string),
				})
			}
		}

		// Si la recherche par ingrédient est fructueuse, récupérez les cocktails correspondants
		if searchByIngredientData["drinks"] != nil {
			for _, drink := range searchByIngredientData["drinks"].([]interface{}) {
				d := drink.(map[string]interface{})
				cocktails = append(cocktails, Cocktail{
					Name:           d["strDrink"].(string),
					Instructions:   d["strInstructions"].(string),
					ThumbnailImage: d["strDrinkThumb"].(string),
				})
			}
		}

		// Affichez les résultats de la recherche
		t, err := template.ParseFiles("templates/results.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, cocktails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
