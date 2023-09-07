package main

import (
	// "errors"
	// "log"
	// "fmt"
	// "database/sql"
	"fmt"
	"html/template"
	"net/http"
	"plantyplantman/go-htmx-server/database"
	"strconv"

	// "plantyplantman/go-htmx-server/database"
	// "strconv"
	"time"

	// "html/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ProductTableHeader struct {
	Header string
}

type ProductTable struct {
	ProductTableHeaders []ProductTableHeader
	Products            []database.Product
}

var PRODUCT_TABLE_HEADERS = []ProductTableHeader{
				{"Sku"},
				{"Product Name"},
				{"Price"},
				{"Promo Price"},
				{"Soh"},
			}


func main() {
	db, err := database.Connect()
	if err != nil {
		fmt.Printf("%v\n\n", err)
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Println("err")
		}

		products, err := database.GetAllProducts(db, 55, 80)
		if err != nil {
			executeErrorTemplate(w, TemplateErrorMessage{Message: "database.GetAllProducts err"}, "database.GetAllProducts err")
		}

		data := ProductTable{
			ProductTableHeaders: PRODUCT_TABLE_HEADERS,
			Products: products,
		}
		err = tpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			fmt.Printf("%v\n", err)
			executeErrorTemplate(w, TemplateErrorMessage{fmt.Sprintf("Function: tpl.ExecuteTemplate(index.html)<br>Error: %v", err)},
				"executeProductTable")
		}
	})

	r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		var err error
		err = r.ParseForm()
		if err != nil {
			executeErrorTemplate(w, TemplateErrorMessage{fmt.Sprintf("Function: r.ParseForm()<br>Error: %v", err)},"parseForm")
		}

		if sku := r.Form["SkuField"]; len(sku) > 0 && sku[0] != "" {
			intSku, err := strconv.ParseInt(sku[0], 10, 64)
			if err != nil {
				executeErrorTemplate(w, TemplateErrorMessage{fmt.Sprintf("Function: strconv.ParseUint(sku, 10, 64)<br>Error: %v", err)},"parseForm")
			}
			product, err := database.GetProductFromSku(db, intSku)
			if err != nil {
				executeErrorTemplate(w, TemplateErrorMessage{fmt.Sprintf("Function: database.GetProductFromSku(db, intSku)<br>Error: %v", err)},"parseForm")
			}

			data := ProductTable{
				ProductTableHeaders: PRODUCT_TABLE_HEADERS,
				Products: []database.Product{product},
			}

			err = executeProductTable(w, data)
			if err != nil {
				executeErrorTemplate(w, TemplateErrorMessage{fmt.Sprintf("Function: executeProductTable called from /products<br>Error: %v", err)},"parseForm")
			}
		}
		// prodName := r.Form["prodName"]
	})


	http.ListenAndServe(":42069", r)
}



func executeProductTable(w http.ResponseWriter, data ProductTable) error {
	tmplStr := `<table id="ProductTable">
      <tr id="ProductTableHeaders">
        {{range .ProductTableHeaders}}
        <th>{{.Header}}</th>
        {{end}}
      </tr>

      <div id="Products">
        {{range .Products}}
        <tr id="{{.Sku}}">
          <td>{{.Sku}}</td>
          <td>{{.ProdName}}</td>
          <td>{{.Soh}}</td>
          <td>{{.Price}}</td>
          <td>{{.PromoPrice}}</td>
        </tr>
        {{end}}
      </div>
    </table>`

	tmpl, err := template.New("ProductTable").Parse(tmplStr)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "ProductTable", data)
}

type TemplateErrorMessage struct {
	Message string
}

func executeErrorTemplate(w http.ResponseWriter, errorMessage TemplateErrorMessage, name string) error {
	tmplStr := `<div id="ErrorMessage">UHOH<br/>SOMETHING WENT WRONG!<br/>{{.Message}}`
	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, name, errorMessage)
}

// RESTy routes for "articles" resource

// r.Route("/articles", func(r chi.Router) {
//   r.With(paginate).Get("/", listArticles)                           // GET /articles
//   r.With(paginate).Get("/{month}-{day}-{year}", listArticlesByDate) // GET /articles/01-16-2017
//
//   r.Post("/", createArticle)                                        // POST /articles
//   r.Get("/search", searchArticles)                                  // GET /articles/search
//
//   // Regexp url parameters:
//   r.Get("/{articleSlug:[a-z-]+}", getArticleBySlug)                // GET /articles/home-is-toronto
//
//   // Subrouters:
//   r.Route("/{articleID}", func(r chi.Router) {
//     r.Use(ArticleCtx)
//     r.Get("/", getArticle)                                          // GET /articles/123
//     r.Put("/", updateArticle)                                       // PUT /articles/123
//     r.Delete("/", deleteArticle)                                    // DELETE /articles/123
//   })
// })
//
// // Mount the admin sub-router
// r.Mount("/admin", adminRouter())
