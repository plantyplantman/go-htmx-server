package parsers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"plantyplantman/go-htmx-server/database"
	"slices"
	"strconv"
	"strings"
	"time"
)

type StoreStockReportLine struct {
	Sku         *string
	ProdName    *string
	Soh         *int
	UnitCost    *float64
	TotalCost   *float64
	LastOrdered *time.Time
}

type StoreStockReport struct {
	Store database.Store
	Lines []StoreStockReportLine
}

type CombinedStockReport struct {
	Sku         *uint64
	ProdName    *string
	SohPetrie   *int
	SohFranklin *int
	SohBunda    *int
	SohCon      *int
	SohCombined *int
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

func skipLines(reader *csv.Reader, lines int) error {
	for i := 0; i < lines; i++ {
		_, err := reader.Read()
		if err != nil {
			return err
		}
	}
	return nil
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

	prodNameStr := strings.TrimSpace(prodName)
	sohInt, err := strconv.Atoi(soh)
	if err != nil {
		log.Printf("Error converting soh to int: %s\nProduct: %v\nUsing default value 0", err, line)
		sohInt = 0
	}
	unitCostFloat, err := parseFloat(unitCost)
	if err != nil {
		log.Printf("Error converting unitCost to float: %s\nProduct: %v\nUsing default value 0", err, line)
	}
	totalCostFloat, err := parseFloat(totalCost)
	if err != nil {
		log.Printf("Error converting totalCost to float: %s\nProduct: %v\nUsing default value 0", err, line)
	}
	lastOrderedTime, err := time.Parse("2/1/06", lastOrdered)
	if err != nil {
		log.Printf("Error converting lastOrdered to time.Time: %s\nProduct: %v\nUsing default value 0", err, line)
	}
	return StoreStockReportLine{
		Sku:         &sku,
		ProdName:    &prodNameStr,
		Soh:         &sohInt,
		UnitCost:    &unitCostFloat,
		TotalCost:   &totalCostFloat,
		LastOrdered: &lastOrderedTime,
	}, nil
}

type ProductRetailListLine struct {
	Mnpn     *uint64
	Sku      *uint64
	ProdNo   *uint64
	ProdName *string
	Price    *float64
}

type ProductRetailList struct {
	Lines []ProductRetailListLine
}

// func productRetailListParser(line []string) (ProductRetailListLine, error) {
// 	relevantCols := []int{1, 3, 5, 6, 7, 8}
//
// 	relevantFields := getRelevantFields(line, relevantCols)
// 	mnpn := relevantFields[0]
// 	sku := relevantFields[1]
// 	prodNo := relevantFields[2]
// 	prodName := relevantFields[3]
// 	price := relevantFields[4]
//
// 	// mnpn
// 	mnpnInt, err := strconv.ParseUint(mnpn, 10, 64)
// 	if err != nil {
// 		log.Printf("Error converting mnpn to int: %s\nProduct: %v", err, line)
// 		mnpnInt = 0
// 	}
//
// }

func ParseStockReport(file string, store database.Store) (StoreStockReport, error) {
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
		Store: store,
		Lines: products,
	}, nil
}
