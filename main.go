// package main
//
// import (
// 	"html/template"
// 	"net/http"
// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"log"
// 	"plantyplantman/go-htmx-server/database"
//   "plantyplantman/go-htmx-server/parsers"
// )
// func upsertProducts(report parsers.StoreStockReport) {
// 	db := database.Connect()
// 	defer db.Close()
// 	const zero float64 = 0
//
// 	for _, line := range report.Lines {
// 		product := database.Product{
// 			Sku:      *line.Sku,
// 			ProdName: *line.ProdName,
//       Price: 0,
//       PromoPrice: 0,
// 		}
// 		data, err := database.UpsertProduct(db, product)
// 		if err != nil {
// 			log.Printf("Error upserting product: %s\nProduct: %v", err, line)
// 			continue
// 		}
// 		log.Printf("Upserted product: %v", data)
// 	}
//
// }
//
// func getProducts(page int, pageLimit int) ([]database.Product, error) {
// 	db := database.Connect()
// 	defer db.Close()
//
// 	products, err := database.GetAllProducts(db, page, pageLimit)
// 	if err != nil {
// 		log.Printf("Error getting all products: %s", err)
// 		return nil, err
// 	}
// 	return products, nil
// }
//
//
// func main() {
//   db := database.Connect()
// 	r := chi.NewRouter()
//
//   // Middleware
// 	r.Use(middleware.Logger)
//
//
// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 	   tmpl, _ := template.ParseFiles("index.html")
//      product, err := database.GetProductFromSku(db, 9300711776258)
//      if err != nil {
//        product = database.Product{}
//      }
// 	   data := map[string]database.Product{"product": product}
// 	   tmpl.ExecuteTemplate(w, "index.html", data)
// 	})
//
// //   r.Get("/products/1", func( w http.ResponseWriter, r *http.Request) {
// //     tmplStr := `
// //     <table>
// //   <tr>
// //     <th>Sku</th>
// //     <th>Product Name</th>
// //     <th>Price</th>
// //     <th>Promo Price</th>
// //     <th>SOH Petrie</th>
// //     <th>SOH Manuka</th>
// //     <th>SOH Bunda</th>
// //     <th>SOH Con</th>
// //     <th>SOH Total</th>
// //   </tr>
// //   <tr>
// //     <td id=\"sku\">{{.Sku}}</td>
// //     <td id="\prodName\">{{.ProdName}}</td>
// //     <td id="\price\">{{.Price}}</td>
// //     <td id="\promoPrice\">{{.PromoPrice}}</td>
// //     <td id="\sohPetrie\">{{.SohPetrie}}</td>
// //     <td id="\sohManuka\">{{.SohManuka}}</td>
// //     <td id="\sohBunda\">{{.SohBunda}}</td>
// //     <td id="\sohCon\">{{.SohCon}}</td>
// //     <td id="\sohTotal\">{{.SohTotal}}</td>
// //   </tr>
// // </table>
// //     `
// //     product, err := database.GetProductFromSku(db, 9300711776258)
// //     if err != nil {
// //       product = database.Product{}
// //     }
// //   })
//
//   r.Get("/stores/1", func(w http.ResponseWriter, r *http.Request) {
//     tmplStr := `
//     <div hx-target="this" hx-swap="outerHTML">
//       <div><label>Store Name</label>: Petrie</div>
//       <button hx-get="/stores/1/edit" class="btn btn-primary">
//         Click To Edit
//       </button>
//     </div>`
//     tmpl, _ := template.New("storeEditingForm").Parse(tmplStr)
//
//     tmpl.ExecuteTemplate(w, "storeEditingForm", map[string]int{})
//   })
//
//   r.Get("/stores/1/edit", func(w http.ResponseWriter, r *http.Request) {
//     tmplStr := `
// <form hx-put="/stores/1" hx-target="this" hx-swap="outerHTML">
//   <div>
//     <label>Store</label>
//     <input type="text" name="Store" value="PETRIE" />
//   </div>
//   <button class="btn">Submit</button>
//   <button class="btn" hx-get="/stores/1">Cancel</button>
// </form>
// `
//   tmpl, _ := template.New("storeEditForm").Parse(tmplStr)
//   tmpl.ExecuteTemplate(w, "storeEditForm", map[string]int{})
//   })
//
//
//
// 	//  r.Post("/products", func(w http.ResponseWriter, _ *http.Request) {
// 	// 	tmplStr := "<div id=\"product\">{{.product}}</div>"
// 	// 	tmpl := template.Must(template.New("product").Parse(tmplStr))
// 	// 	counter.Increase()
// 	// 	data := map[string]int{
// 	// 		"CounterValue": counter.GetValue(),
// 	// 	}
// 	// 	tmpl.ExecuteTemplate(w, "counter", data)
// 	// })
// 	// r.Post("/decrease", func(w http.ResponseWriter, _ *http.Request) {
// 	// 	tmplStr := "<div id=\"counter\">{{.CounterValue}}</div>"
// 	// 	tmpl := template.Must(template.New("counter").Parse(tmplStr))
// 	// 	counter.Decrease()
// 	// 	data := map[string]int{
// 	// 		"CounterValue": counter.GetValue(),
// 	// 	}
// 	// 	tmpl.ExecuteTemplate(w, "counter", data)
// 	// })
// 	//
//
// 	http.ListenAndServe(":3000", r)
// }

package main

import (
	"errors"
	"log"
	"net/http"
	"plantyplantman/go-htmx-server/database"
	"strconv"
	"time"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	hxhttp "github.com/maragudk/gomponents-htmx/http"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
	ghttp "github.com/maragudk/gomponents/http"
)

func main() {
	if err := start(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func start() error {
	now := time.Now()
	mux := http.NewServeMux()
  db, err := database.Connect()
  if err != nil {
    log.Panicln("Connection to database failed")
  }
  log.Println("Connection to database succeeded")

	mux.HandleFunc("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		if r.Method == http.MethodPost && hxhttp.IsBoosted(r.Header) {

      products, err := database.GetAllProducts(db, 1, 20)
      if err != nil {
        products = []database.Product{}
      }

			hxhttp.SetPushURL(w.Header(), "/?time="+now.Format(timeFormat))

			return productTable(products), nil
		}
		return root(now), nil
	}))

  mux.HandleFunc("/products/sku", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
    if r.Method == http.MethodPost && hxhttp.IsBoosted(r.Header) {
      err := r.ParseForm()
      if err != nil {
        log.Panicln("PARSEFORM FAILED IN POST /products")
      }

      sku := r.Form["sku-input"][0]
      hxhttp.SetPushURL(w.Header(), "/products/sku?="+sku)

      intSku, err := strconv.ParseInt(sku, 10, 64)
      if err != nil {
        intSku = 0
      }

      if intSku == 0 {
        data, err := database.GetAllProducts(db, 1, 20)
        if err != nil {
          return productTable([]database.Product{}), nil
        }
        return productTable(data),nil
      }

      product, err := database.GetProductFromSku(db, intSku)

      if err != nil {
        product = database.Product{}
      }

      r := []database.Product{product}
      return productTable(r), nil
    }
    return nil,nil
  }))

	log.Println("Starting on http://localhost:8080")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

const timeFormat = "15:04:05"

func root(now time.Time) g.Node {
	return c.HTML5(c.HTML5Props{
		Title: now.Format(timeFormat),
		Head: []g.Node{
			Script(Src("https://cdn.tailwindcss.com?plugins=forms,typography")),
			Script(Src("https://unpkg.com/htmx.org")),
		},
		Body: []g.Node{
			Div(Class("max-w-7xl mx-auto p-4 prose lg:prose-lg xl:prose-xl"),
				H1(g.Text(`PRODUCTS`)),
				FormEl(Method("post"), Action("/products/sku"), hx.Boost("true"), hx.Target("#productTableHeaders"), hx.Swap("outerHTML"),
          Label(For("Sku"), g.Text("SKU")),
          Input(Type("text"), ID("skuInput"), Name("sku-input")),
					Button(Type("submit"), g.Text(`Get Products`),
						Class("rounded-md border border-transparent bg-orange-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2"),
					),
				),
        productTableHeaders(true),
			),
		},
	})
}

func partial(now time.Time) g.Node {
	return P(ID("partial"), g.Textf(`Time was last updated at %v.`, now.Format(timeFormat)))
}

func productItem(product database.Product) g.Node {
  return Tr(Td(g.Textf("%v",product.Sku)),
            Td(g.Textf("%v", product.ProdName)),
            Td(g.Textf("%v", product.Price)),
            Td(g.Textf("%v", product.PromoPrice)),
          )
}

func productTableHeaders(table bool) g.Node {
  tableHeaders := Tr(ID("productTableHeaders"),
                     Th(g.Text("Sku")),
                     Th(g.Text("Name")),
                     Th(g.Text("Price")),
                     Th(g.Text("PromoPrice")),
                     Th(g.Text("Petrie SOH")),
                     Th(g.Text("Franklin SOH")),
                     Th(g.Text("Bunda SOH")),
                     Th(g.Text("Con SOH")),
                     Th(g.Text("Total SOH")),
                   )
  if table {
    return Table(tableHeaders)
  }
  return tableHeaders
}

func productTable(products []database.Product) g.Node {
  tableRows := g.Group(g.Map(products, func(product database.Product) g.Node {
    return productItem(product)
  }))

  return Table(ID("productTable"), productTableHeaders(false), tableRows)
}
