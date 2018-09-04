package main

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/robfig/cron"
	"golang.org/x/text/encoding/charmap"
)

const (
	pplUpdateURL = "https://www.deutschepost.de/content/dam/mlm.nf/dpag/technische_downloads/update_internetmarke/ppl_update.xml"
)

var productState *ProductsState

type ProductsState struct {
	version    string
	products   []Product
	productMap map[int]Product
	countries  []Country
	countryMap map[string]Country //Shortcode -> country
	sync.RWMutex
}

type Product struct {
	// Parsing this undocumented mess is still better than using yet another
	// SOAP API that spams you 3 times a day via e-mail
	Tracking      bool   // Col. 1
	ProductId     int    // Col. 2
	International bool   // Col. 3
	Name          string // Col. 4
	Cost          int    // Col 5, given in euro but stored in cent internally
	// Col. 13 "Zustellkennung", A: Age check, R: Registered, H: High Priority, w: Product shipment
	// Col. 13 == A provides a different ordering, skip those for now
	// Col. 14: Minimal length, mm
	// Col. 15: Minimal width, mm
	// Col. 16: Minimal height, mm
	// Col. 17: Maximal length, mm
	// Col. 18: Maximal width, mm
	// Col. 19: Maximal height, mm
	// Col. 20: Minimal weight, gramm
	MaxWeight   int    // Col. 21, gramm
	Description string // Col. 45
}

type Country struct {
	Fullname string
	Code     string
}

type PplUpdate struct {
	App []struct {
		Updates struct {
			Li struct {
				Version      string `xml:"version"`
				Service      string `xml:"service"`
				UpdateLink   string `xml:"updateLink"`
				Description  string `xml:"description"`
				Dependencies struct {
					Dependency []struct {
						Name string `xml:"name,attr"`
						URL  string `xml:",chardata"`
					} `xml:"dependency"`
				} `xml:"dependencies"`
			} `xml:"li"`
		} `xml:"updates"`
	} `xml:"app"`
}

//TODO proper error handling, but this has to work anyway
func updateProductState() {

	log.Println("Updating product state")

	if productState != nil {
		productState.Lock()
		defer productState.Unlock()
	}

	state := new(ProductsState)

	ver, link, linkCountry, err := getLatestProductListVersion()
	if err != nil {
		log.Fatal(err)
	}

	if productState != nil && productState.version == ver {
		log.Println("Product list version up-to-date:", ver)
		return
	}

	state.version = ver

	list, err := getProductList(link)
	if err != nil {
		log.Fatal(err)
	}

	state.products = list

	state.productMap = make(map[int]Product)
	for _, v := range list {
		state.productMap[v.ProductId] = v
	}

	listCountry, err := getCountryList(linkCountry)
	if err != nil {
		log.Fatal(err)
	}

	state.countries = listCountry

	state.countryMap = make(map[string]Country)
	for _, v := range listCountry {
		state.countryMap[v.Code] = v
	}

	productState = state

	log.Println("Updated product list to version:", ver)
}

func init() {
	updateProductState()
	c := cron.New()
	err := c.AddFunc("@daily", updateProductState)
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
}

func getLatestProductListVersion() (version string, linkProduct string, linkCountry string, err error) {

	var response *http.Response
	response, err = http.DefaultClient.Get(pplUpdateURL)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New("Error requesting product update xml: " + response.Status)
		return
	}
	defer response.Body.Close()

	var productListXml []byte
	productListXml, err = ioutil.ReadAll(response.Body)

	if err != nil {
		err = errors.New("Error reading product update xml: " + err.Error())
		return
	}

	ppl := new(PplUpdate)
	err = xml.Unmarshal(productListXml, ppl)

	if err != nil {
		return
	}

	version = strings.TrimSpace(ppl.App[0].Updates.Li.Version)
	linkProduct = strings.TrimSpace(ppl.App[0].Updates.Li.UpdateLink)

	for _, v := range ppl.App[0].Updates.Li.Dependencies.Dependency {
		if strings.ToLower(strings.TrimSpace(v.Name)) == "countries" {
			linkCountry = strings.TrimSpace(v.URL)
		}
	}
	return
}

func getCountryList(url string) (countries []Country, err error) {
	// Get countries
	var countryResponse *http.Response
	countryResponse, err = http.DefaultClient.Get(url)
	if err != nil {
		return
	}
	if countryResponse.StatusCode != http.StatusOK {
		err = errors.New("Error requesting product list csv: " + countryResponse.Status)
		return
	}
	defer countryResponse.Body.Close()

	charsetDecorder := charmap.Windows1252.NewDecoder()

	countryReader := csv.NewReader(charsetDecorder.Reader(countryResponse.Body))
	countryReader.Comma = ';'
	countryReader.TrimLeadingSpace = true

	var countryRecords [][]string
	countryRecords, err = countryReader.ReadAll()

	if err != nil {
		err = errors.New("Error reading country list csv: " + err.Error())
		return
	}

	for _, v := range countryRecords {
		var c Country
		c.Fullname = strings.TrimSpace(v[0])
		c.Code = strings.TrimSpace(v[1])
		countries = append(countries, c)
	}

	return
}

func getProductList(url string) (products []Product, err error) {

	var response *http.Response
	response, err = http.DefaultClient.Get(url)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New("Error requesting product list csv: " + response.Status)
		return
	}
	defer response.Body.Close()

	charsetDecorder := charmap.Windows1252.NewDecoder()

	reader := csv.NewReader(charsetDecorder.Reader(response.Body))
	reader.Comma = ';'
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()

	if err != nil {
		err = errors.New("Error reading product list csv: " + err.Error())
		return
	}

	for _, v := range records {
		// The answer to life, the universe and everything
		if len(v) < 42 {
			err = fmt.Errorf("Too few columns in product list: %d", len(v))
			return
		}

		if strings.ToUpper(strings.TrimSpace(v[13])) == "A" {
			//*A*lterssichtprüfung?
			//Skip for now due to different column ordering
			continue
		}

		var p Product

		trackingTxt := strings.TrimSpace(strings.ToLower(v[1]))
		var trackingBool bool
		switch trackingTxt {
		case "ja":
			trackingBool = true
		case "nein":
			trackingBool = false
		default:
			err = fmt.Errorf("Error parsing Ja/Nein bool: '%s'", trackingTxt)
			return
		}

		p.Tracking = trackingBool

		p.ProductId, err = strconv.Atoi(strings.TrimSpace(v[2]))
		if err != nil {
			err = errors.New("Parsing error: " + err.Error())
			return
		}

		p.International, err = strconv.ParseBool(strings.TrimSpace(v[3]))
		if err != nil {
			err = errors.New("Parsing error: " + err.Error())
			return
		}

		p.Name = strings.TrimSpace(v[4])
		p.Name = strings.Replace(p.Name, "Intern.", "", -1)       // Displayed separately
		p.Name = strings.Replace(p.Name, "International", "", -1) // Displayed separately
		p.Name = strings.Replace(p.Name, "Einsch+Eigh", "Einschreiben Eigenhändig", -1)

		var costEuro float64
		costEuro, err = strconv.ParseFloat(strings.Replace(strings.TrimSpace(v[5]), ",", ".", 1), 64)
		if err != nil {
			err = errors.New("Parsing error: " + err.Error())
			return
		}
		costCent := Round(costEuro*100, .5, 0)

		if math.Mod(costCent, 1) != 0 {
			err = fmt.Errorf("Invalid cent cost %v for euro value %s", costCent, v[5])
			return
		}

		p.Cost = int(costCent)

		p.MaxWeight, err = strconv.Atoi(strings.TrimSpace(v[21]))
		if err != nil {
			err = errors.New("Parsing error: " + err.Error())
			return
		}

		p.Description = strings.TrimSpace(v[41])
		products = append(products, p)

	}

	return
}
