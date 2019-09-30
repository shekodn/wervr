package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "os"

  "github.com/auth0-community/auth0"
  jose "gopkg.in/square/go-jose.v2"
  "github.com/joho/godotenv"
  "github.com/gorilla/handlers"
  "github.com/sirupsen/logrus"
  "github.com/gorilla/mux"
)

var log = logrus.New()

func main() {

  e := godotenv.Load()

  if e != nil {
    log.Println("ENV: ", e)
  }

  audience := os.Getenv("AUDIENCE")
  secret := os.Getenv("SECRET")

  if audience == "" {
    log.Fatal("audience is not set.")
  }

  if secret == "" {
    log.Fatal("secret is not set.")
  }

  r := mux.NewRouter()

  r.Handle("/", http.FileServer(http.Dir("./views/")))

  // Our API is going to consist of three routes
  // /status - which we will call to make sure that our API is up and running
  r.Handle("/status", StatusHandler).Methods("GET")

  /* We will add the middleware to our products and feedback routes. The status route will be publicly accessible */

  // /products - which will retrieve a list of products that the user can leave feedback on
  r.Handle("/products", authMiddleware(ProductsHandler)).Methods("GET")

  // /products/{slug}/feedback - which will capture user feedback on product
  r.Handle("/products/{slug}/feedback", authMiddleware(AddFeedbackHandler)).Methods("POST")

  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

  http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))

}

/* HANDLERS */
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
  w.Write([]byte("Not Implemented"))
})

/* The status handler will be invoked when the user calls the /status route
   It will simply return a string with the message "API is up and running" */
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
  w.Write([]byte("API is up and running"))
})


func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        secret := []byte(os.Getenv("SECRET"))
        secretProvider := auth0.NewKeyProvider(secret)

        audience := os.Getenv("AUDIENCE")
        audiences := []string{audience}

        configuration := auth0.NewConfiguration(secretProvider, audiences, "https://" + os.Getenv("AUTH0-DOMAIN") + ".auth0.com/", jose.HS256)

        validator := auth0.NewValidator(configuration, nil)

        token, err := validator.ValidateRequest(r)

        if err != nil {
            fmt.Println(err)
            fmt.Println("Token is not valid:", token)
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(err.Error()))
        } else {
            next.ServeHTTP(w, r)
        }
    })
}


/* The products handler will be called when the user makes a GET request to the /products endpoint.
   This handler will return a list of products available for users to review */
var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
  // Here we are converting the slice of products to JSON
  payload, _ := json.Marshal(products)

  w.Header().Set("Content-Type", "application/json")
  w.Write([]byte(payload))
})

/* The feedback handler will add either positive or negative feedback to the product
   We would normally save this data to the database - but for this demo, we'll fake it
   so that as long as the request is successful and we can match a product to our catalog of products
   we'll return an OK status. */
var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var product Product
    vars := mux.Vars(r)
    slug := vars["slug"]

    for _, p := range products {
        if p.Slug == slug {
            product = p
        }
    }

    w.Header().Set("Content-Type", "application/json")
    if product.Slug != "" {
        payload, _ := json.Marshal(product)
        w.Write([]byte(payload))
    } else {
        w.Write([]byte("Product Not Found"))
    }
})


/* We will first create a new type called Product
   This type will contain information about VR experiences */
type Product struct {
    Id int
    Name string
    Slug string
    Description string
}

/* We will create our catalog of VR experiences and store them in a slice. */
var products = []Product{
  Product{Id: 1, Name: "Hover Shooters", Slug: "hover-shooters", Description : "Shoot your way to the top on 14 different hoverboards"},
  Product{Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer", Description : "Explore the depths of the sea in this one of a kind underwater experience"},
  Product{Id: 3, Name: "Dinosaur Park", Slug : "dinosaur-park", Description : "Go back 65 million years in the past and ride a T-Rex"},
  Product{Id: 4, Name: "Cars VR", Slug : "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
  Product{Id: 5, Name: "Robin Hood", Slug: "robin-hood", Description : "Pick up the bow and arrow and master the art of archery"},
  Product{Id: 6, Name: "Real World VR", Slug: "real-world-vr", Description : "Explore the seven wonders of the world in VR"},
}
