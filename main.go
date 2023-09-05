package main

import (
	// "html/template"
	// "net/http"
	"encoding/csv"
	"fmt"

	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	// "bufio"
	// "fmt"
	"log"
	"os"
	"plantyplantman/go-htmx-server/database"
	"slices"
	"strconv"
	"strings"
	"time"
)

type StoreStockReportLine struct {
	sku         *uint64
	prodName    *string
	soh         *int
	unitCost    *float64
	totalCost   *float64
	lastOrdered *time.Time
}

type StoreStockReport struct {
	store Store
	lines []StoreStockReportLine
}

type CombinedStockReport struct {
	sku         *uint64
	prodName    *string
	sohPetrie   *int
	sohFranklin *int
	sohBunda    *int
	sohCon      *int
	sohCombined *int
}

type Store int64

const (
	PETRIE Store = iota
	BUNDA
	FRANKLIN
	CON
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getRelevantFields(data []string, relevantCols []int) []string {
	fltd := make([]string, 0, len(data))
	for i, e := range data {
		if slices.Contains(relevantCols, i) {
			fltd = append(fltd, e)
		}
	}
	return fltd
}

func removeCommaFromNumber(num string) string {
	return strings.ReplaceAll(num, ",", "")
}

func parseFloat(num string) (float64, error) {
	num = removeCommaFromNumber(num)
	return strconv.ParseFloat(num, 64)
}

func stockReportParser(line []string) (StoreStockReportLine, error) {
	relevantCols := []int{1, 3, 5, 6, 7, 8}

	relevantFields := getRelevantFields(line, relevantCols)
	sku := relevantFields[0]
	prodName := relevantFields[1]
	soh := relevantFields[2]
	unitCost := relevantFields[3]
	totalCost := relevantFields[4]
	lastOrdered := relevantFields[5]

	// sku
	skuInt, err := strconv.ParseUint(sku, 10, 64)
	if err != nil {
		log.Printf("Error converting sku to int: %s\nProduct: %v", err, line)
		return StoreStockReportLine{}, err
	}

	// prodName
	prodNameStr := strings.TrimSpace(prodName)

	// soh
	sohInt, err := strconv.Atoi(soh)
	if err != nil {
		log.Printf("Error converting soh to int: %s\nProduct: %v\nUsing default value 0", err, line)
		sohInt = 0
	}

	// unitCost
	unitCostFloat, err := parseFloat(unitCost)
	if err != nil {
		log.Printf("Error converting unitCost to float: %s\nProduct: %v\nUsing default value 0", err, line)
	}

	// totalCost
	totalCostFloat, err := parseFloat(totalCost)
	if err != nil {
		log.Printf("Error converting totalCost to float: %s\nProduct: %v\nUsing default value 0", err, line)
	}

	// lastOrdered
	lastOrderedTime, err := time.Parse("2/1/06", lastOrdered)
	if err != nil {
		log.Printf("Error converting lastOrdered to time.Time: %s\nProduct: %v\nUsing default value 0", err, line)
	}

	return StoreStockReportLine{
		sku:         &skuInt,
		prodName:    &prodNameStr,
		soh:         &sohInt,
		unitCost:    &unitCostFloat,
		totalCost:   &totalCostFloat,
		lastOrdered: &lastOrderedTime,
	}, nil
}

func skipLines(reader *csv.Reader, lines int) error {
	for i := 0; i < lines; i++ {
		_, err := reader.Read()
		if err != nil {
			return err
		}
	}
	return nil
}

func parseStockReport(file string, store Store) (StoreStockReport, error) {
	f, err := os.Open(file)
	if err != nil {
		return StoreStockReport{}, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	csvr.Comma = '\t'

	products := make([]StoreStockReportLine, 0, 1000)

	// Read and discard the first and second line (header)
	err = skipLines(csvr, 2)
	if err != nil {
		return StoreStockReport{}, err
	}

	// Read the rest of the lines

	for {
		line, err := csvr.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return StoreStockReport{}, err
		}
		product, err := stockReportParser(line)
		if err != nil {
			log.Printf("Error parsing product: %s\nProduct: %v", err, line)
			continue
		}
		products = append(products, product)
	}

	fmt.Println("Finished parsing stock report")

	return StoreStockReport{
		store: store,
		lines: products,
	}, nil
}

func upsertProducts(report StoreStockReport) {
	db := database.Connect()
	defer db.Close()
	const zero float64 = 0

	for _, line := range report.lines {
		product := database.Product{
			Sku:      *line.sku,
			ProdName: *line.prodName,
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
	//  counter := &Counter{}
	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//    tmpl, _ := template.ParseFiles("index.html")
	//    data := map[string]int{
	//      "CounterValue": counter.GetValue(),
	//    }
	//    tmpl.ExecuteTemplate(w, "index.html", data)
	// })
	//  r.Post("/increase", func(w http.ResponseWriter, _ *http.Request) {
	// 	tmplStr := "<div id=\"counter\">{{.CounterValue}}</div>"
	// 	tmpl := template.Must(template.New("counter").Parse(tmplStr))
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
	// http.ListenAndServe(":3000", r)

	// path := "/Users/home/Desktop/Work/Stock reports/230904/BUNDA.txt"
	// lines, _ := parseStockReport(path, BUNDA)
	// if err != nil {
	// 	log.Fatalf("readLines: %s", err)
	// }

	// for _, line := range lines.lines {
	// 	fmt.Println(line)
	// }

	// data, err := parseStockReport("/Users/home/Desktop/Work/Stock reports/230904/BUNDA.txt", BUNDA)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// upsertProducts(data)

	products, err := getProducts(100,10)
	if err != nil {
		log.Fatal(err)
	}

	for _, product := range products {
		log.Printf("%v", product)
	}
}
