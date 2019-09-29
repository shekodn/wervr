package main

import (
  "encoding/json"
  "net/http"
  "os"
  "time"

  "github.com/auth0/go-jwt-middleware"
  "github.com/gorilla/handlers"
  jwt "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
)

func main() {
  r := mux.NewRouter()

  r.Handle("/", http.FileServer(http.Dir("./views/")))

  // Our API is going to consist of three routes
  // /status - which we will call to make sure that our API is up and running
  r.Handle("/status", StatusHandler).Methods("GET")

  /* We will add the middleware to our products and feedback routes. The status route will be publicly accessible */

  // /products - which will retrieve a list of products that the user can leave feedback on
  r.Handle("/products", jwtMiddleware.Handler(ProductsHandler)).Methods("GET")

  // /products/{slug}/feedback - which will capture user feedback on product
  r.Handle("/products/{slug}/feedback", jwtMiddleware.Handler(AddFeedbackHandler)).Methods("POST")

  // AUTH
  r.Handle("/get-token", GetTokenHandler).Methods("GET")


  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

  http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))

}

/* Set up a global string for our secret */
var mySigningKey = []byte("secret")

/* HANDLERS */
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
  w.Write([]byte("Not Implemented"))
})

/* The status handler will be invoked when the user calls the /status route
   It will simply return a string with the message "API is up and running" */
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
  w.Write([]byte("API is up and running"))
})

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
var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
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

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    /* Create the token */
  token := jwt.New(jwt.SigningMethodHS256)

  /* Create a map to store our claims */
  claims := token.Claims.(jwt.MapClaims)

  /* Set token claims */
  claims["admin"] = true
  claims["name"] = "Ado Kukic"
  claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

  /* Sign the token with our secret */
  tokenString, _ := token.SignedString(mySigningKey)

  /* Finally, write the token to the browser window */
  w.Write([]byte(tokenString))
})

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
  ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
    return mySigningKey, nil
  },
  SigningMethod: jwt.SigningMethodHS256,
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
