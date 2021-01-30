package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var Articles []*Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "¡Bienvenido a la página de inicio!")
	fmt.Println("Llegada al endpoint: homePage")
}

func handleRequest() {
	/*	//route to homePage
		http.HandleFunc("/", homePage)
		//route to returnAllArticles
		http.HandleFunc("/articles", returnAllArticles)


		log.Fatal(http.ListenAndServe(":10000", nil))*/

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	// NOTE: Ordering is important here! This has to be defined before
	// the other `/article` endpoint.
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Endpoint: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var founded bool = false
	//fmt.Fprintf(w, "Key: "+key)
	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	for _, article := range Articles {
		if article.Id == key {
			founded = true
			json.NewEncoder(w).Encode(article)
		}
	}

	if !founded {
		fmt.Fprintf(w, "Key no encontrada: "+key)
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	//fmt.Println("Ingreso crear nuevo art")
	reqBody, _ := ioutil.ReadAll(r.Body)
	//fmt.Fprintf(w, "%+v", string(reqBody))
	var article *Article
	json.Unmarshal(reqBody, &article)
	// update our global Articles array to include
	// our new Article
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our PUT request
	reqBody, _ := ioutil.ReadAll(r.Body)

	var articleToUpdate Article
	json.Unmarshal(reqBody, &articleToUpdate)
	// update our article array with new values
	for _, article := range Articles {
		if article.Id == articleToUpdate.Id {
			article.Title = articleToUpdate.Title
			article.Desc = articleToUpdate.Desc
			article.Content = articleToUpdate.Content

			json.NewEncoder(w).Encode(article)
		}
	}
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
	// wish to delete
	id := vars["id"]

	// we then need to loop through all our articles
	for index, article := range Articles {
		// if our id path parameter matches one of our
		// articles
		if article.Id == id {
			// updates our Articles array to remove the
			// article
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}

}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Articles = []*Article{
		{Id: "1", Title: "Hola", Desc: "Descripcion del articulo", Content: "Contenido del articulo"},
		{Id: "2", Title: "Hola2", Desc: "Descripcion del articulo2", Content: "Contenido del articulo2"},
	}

	handleRequest()
}
