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

			return ProductTable(products), nil
		}
		return root(now), nil
	}))


	mux.HandleFunc("/products/search", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (g.Node, error) {
		if r.Method == http.MethodPost && hxhttp.IsBoosted(r.Header) {
			err := r.ParseForm()

			if err != nil {
				log.Panicln("PARSEFORM FAILED IN POST /products")
				return nil, err
			}

			input := r.Form["search-input"][0]
			hxhttp.SetPushURL(w.Header(), "/products/search?="+input)

			products, err := database.Search(db, "Product", "prodName", input)

			if err != nil {
				log.Println(err)
				return nil, err
			}

			return ProductTable(products), nil
		}
			return nil, nil
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
					return ProductTable([]database.Product{}), nil
				}
				return ProductTable(data), nil
			}

			product, err := database.GetProductFromSku(db, intSku)
			if err != nil {
				product = database.Product{}
			}

			retv := []database.Product{product}
			return ProductTable(retv), nil
		}
		return nil, nil
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
			H1(g.Text(`PRODUCTS`)),
			Div(
				Class("max-w-7xl mx-auto p-4 prose lg:prose-lg xl:prose-xl"),
				FormEl(
					Method("post"),
					Action("/products/sku"),
					hx.Boost("true"),
					hx.Target("#productTable"),
					hx.Swap("outerHTML"),
					Label(For("Sku"), g.Text("SKU")),
					Input(Type("text"), ID("skuInput"), Name("sku-input")),
					Button(
						Type("submit"),
						g.Text(`Get Products`),
						Class(
							"rounded-md border border-transparent bg-orange-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2",
						),
					),
				),
			),

			Div(
				Class("max-w-7xl mx-auto p-4 prose lg:prose-lg xl:prose-xl"),
				FormEl(
					Method("post"),
					Action("/products/search"),
					hx.Boost("true"),
					hx.Target("#productTable"),
					hx.Swap("outerHTML"),
					Label(For("search"), g.Text("search")),
					Input(Type("text"), ID("search"), Name("search-input")),
					Button(
						Type("submit"),
						g.Text(`search`),
						Class(
							"rounded-md border border-transparent bg-orange-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:ring-offset-2",
						),
					),
				),
				ProductTable([]database.Product{}),
			),
		},
	})
}

func partial(now time.Time) g.Node {
	return P(ID("partial"), g.Textf(`Time was last updated at %v.`, now.Format(timeFormat)))
}

func productItem(product database.Product) g.Node {
	return Tr(Td(g.Textf("%v", product.Sku)),
		Td(g.Textf("%v", product.ProdName)),
		Td(g.Textf("%v", product.Price)),
		Td(g.Textf("%v", product.PromoPrice)),
	)
}

func ProductTableHeaders(isTable bool) g.Node {
	tableHeadersSlice := []string{
		"Sku",
		"Price",
		"PromoPrice",
		"Petrie SOH",
		"Franklin SOH",
		"Bunda SOH",
		"Con SOH",
		"Total SOH",
	}

	ProductTableHeaders := g.Group(g.Map(tableHeadersSlice, func(s string) g.Node {
		return Th(g.Text(s))
	}))

	if isTable {
		return Table(ProductTableHeaders,ID("productTable"))
	}
	return ProductTableHeaders
}

func ProductTable(products []database.Product) g.Node {

	tableRows := g.Group(g.Map(products, func(product database.Product) g.Node {
		return productItem(product)
	}))

	return Table(ID("productTable"), ProductTableHeaders(false), tableRows)
}
