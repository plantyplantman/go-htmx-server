package main

import (
	"html/template"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"plantyplantman/go-htmx-server/database"
  "plantyplantman/go-htmx-server/parsers"
)
func upsertProducts(report parsers.StoreStockReport) {
	db := database.Connect()
	defer db.Close()
	const zero float64 = 0

	for _, line := range report.Lines {
		product := database.Product{
			Sku:      *line.Sku,
			ProdName: *line.ProdName,
      Price: 0,
      PromoPrice: 0,
		}
		data, err := database.UpsertProduct(db, product)
		if err != nil {
			log.Printf("Error upserting product: %s\nProduct: %v", err, line)
			continue
		}
		log.Printf("Upserted product: %v", data)
	}

}

func getProducts(page int, pageLimit int) ([]database.Product, error) {
	db := database.Connect()
	defer db.Close()

	products, err := database.GetAllProducts(db, page, pageLimit)
	if err != nil {
		log.Printf("Error getting all products: %s", err)
		return nil, err
	}
	return products, nil
}


func main() {
  db := database.Connect()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	   tmpl, _ := template.ParseFiles("index.html")
     product, err := database.GetProductFromSku(db, 9300711776258)
     if err != nil {
       product = database.Product{}
     }
	   data := map[string]database.Product{"product": product}
	   tmpl.ExecuteTemplate(w, "index.html", data)
	})
	//  r.Post("/products", func(w http.ResponseWriter, _ *http.Request) {
	// 	tmplStr := "<div id=\"product\">{{.product}}</div>"
	// 	tmpl := template.Must(template.New("product").Parse(tmplStr))
	// 	counter.Increase()
	// 	data := map[string]int{
	// 		"CounterValue": counter.GetValue(),
	// 	}
	// 	tmpl.ExecuteTemplate(w, "counter", data)
	// })
	// r.Post("/decrease", func(w http.ResponseWriter, _ *http.Request) {
	// 	tmplStr := "<div id=\"counter\">{{.CounterValue}}</div>"
	// 	tmpl := template.Must(template.New("counter").Parse(tmplStr))
	// 	counter.Decrease()
	// 	data := map[string]int{
	// 		"CounterValue": counter.GetValue(),
	// 	}
	// 	tmpl.ExecuteTemplate(w, "counter", data)
	// })
	//
	http.ListenAndServe(":3000", r)
}
