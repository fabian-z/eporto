package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fabian-z/eporto/epservice"
	"html/template"
	"net/url"
	"strings"

	"github.com/gorilla/schema"

	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

const (
	endpoint = "https://internetmarke.deutschepost.de/OneClickForAppV3"
)

//TODO configuration object instead of global state
var (
	Listener              string
	PrinterName           string
	PageFormat            int
	FormatSupportsAddress bool
	PartnerId             string
	KeyPhase              string
	SignatureKey          string
	WalletUser            string
	WalletPassword        string
)

type StampData struct {
	StampAddress  bool
	SenderCompany string
	SenderName    string
	SenderStreet  string
	SenderHouseNo string
	SenderZIP     string
	SenderCity    string
	//SenderCountry hardcoded to "DEU"
	ReceiverCompany string
	ReceiverName    string
	ReceiverStreet  string
	ReceiverHouseNo string
	ReceiverZIP     string
	ReceiverCity    string
	ReceiverCountry string
	Product         int
}

func main() {

	config()

	service, err := epservice.NewOneClickForAppPortTypeV3(endpoint, true, &epservice.AuthenticationData{
		PartnerId: PartnerId,
		KeyPhase:  KeyPhase,
		Key:       SignatureKey,
	})

	if err != nil {
		log.Fatal(err)
	}

	FormatSupportsAddress, err = checkPageFormat(service, PageFormat)

	if err != nil {
		log.Fatal(err)
	}

	h := new(Handler)
	h.service = service

	funcMap := map[string]interface{}{
		"centToEuro": func(cent int) string {
			return fmt.Sprintf("%.2f", Round(float64(cent)/100, .5, 2))
		},

		"newline": func(count int) template.HTMLAttr {

			var s string
			for x := 0; x <= count; x++ {
				s = s + "\n"
			}

			return template.HTMLAttr(s)
		},
	}

	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	h.template = t

	ts, err := template.New("success.html").Funcs(funcMap).ParseFiles("templates/success.html")
	if err != nil {
		log.Fatal(err)
	}
	h.templateSuccess = ts
	h.decoder = schema.NewDecoder()

	log.Println("Starting to listen on", Listener)
	log.Fatal(http.ListenAndServe(Listener, h))
	return

}

func getStampPDFBytes(link string) ([]byte, error) {

	pdfLink, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	if pdfLink.Scheme != "https" {
		return nil, errors.New("Stamp link returned by API is not https, got: " + pdfLink.Scheme)
	}

	var response *http.Response
	response, err = http.DefaultClient.Get(link)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error requesting stamp PDF from URL: %s. Status: %s ", link, response.Status)
	}
	defer response.Body.Close()

	pdf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(pdf[:5], []byte("%PDF-")) {
		return nil, errors.New("Downloaded stamp file is not PDF")
	}

	return pdf, nil
}

type Handler struct {
	service         *epservice.OneClickForAppPortTypeV3
	template        *template.Template
	templateSuccess *template.Template
	decoder         *schema.Decoder
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		productState.RLock()
		defer productState.RUnlock()

		data := struct {
			SupportAddress bool
			Products       []Product
			Countries      []Country
			WalletUser     string
		}{
			FormatSupportsAddress,
			productState.products,
			productState.countries,
			WalletUser,
		}

		err := h.template.Execute(w, data)
		if err != nil {
			log.Println(err)
			return
		}

	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		var stampData StampData
		err = h.decoder.Decode(&stampData, r.PostForm)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		balance, link, err := h.buyAndPrintStamp(&stampData)
		if err != nil && len(link) == 0 {
			http.Error(w, "Error buying stamps: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			WalletUser   string
			Balance      int
			Link         string
			PrintSuccess bool
		}{
			WalletUser,
			balance,
			link,
			err == nil,
		}

		if err := h.templateSuccess.Execute(w, data); err != nil {
			log.Println(err)
		}

	} else {
		http.Error(w, "Method unsupported", http.StatusMethodNotAllowed)
	}

	return
}

func (h *Handler) buyAndPrintStamp(stamp *StampData) (int, string, error) {
	authResponse, err := h.service.AuthenticateUser(&epservice.AuthenticateUserRequest{
		Username: WalletUser,
		Password: WalletPassword,
	})

	if err != nil {
		return 0, "", err
	}

	var pageFormat int32 = int32(PageFormat)
	var productCode int32 = int32(stamp.Product)

	productState.RLock()
	product, ok := productState.productMap[stamp.Product]
	productState.RUnlock()

	if !ok {
		return 0, "", errors.New("Invalid product id")
	}

	if stamp.StampAddress && stamp.ReceiverCountry != "DEU" && !product.International {
		return 0, "", errors.New("Product must be international if receiver is not DEU")
	}

	var total int32 = int32(product.Cost) //Only support one purchase at a time right now
	var createManifest bool = false

	var voucherLayout *epservice.VoucherLayout
	var addressBinding *epservice.AddressBinding

	if FormatSupportsAddress && stamp.StampAddress {
		voucherLayout = &epservice.VoucherLayoutAddressZone

		senderCompany, receiverCompany, senderPerson, receiverPerson := processNames(stamp)
		addressBinding = &epservice.AddressBinding{
			Sender: &epservice.NamedAddress{
				Name: &epservice.Name{
					CompanyName: senderCompany,
					PersonName:  senderPerson,
				},
				Address: &epservice.Address{
					Street:  stamp.SenderStreet,
					HouseNo: stamp.SenderHouseNo,
					Zip:     stamp.SenderZIP,
					City:    stamp.SenderCity,
					Country: "DEU",
				},
			},
			Receiver: &epservice.NamedAddress{
				Name: &epservice.Name{
					CompanyName: receiverCompany,
					PersonName:  receiverPerson,
				},
				Address: &epservice.Address{
					Street:  stamp.ReceiverStreet,
					HouseNo: stamp.ReceiverHouseNo,
					Zip:     stamp.ReceiverZIP,
					City:    stamp.ReceiverCity,
					Country: stamp.ReceiverCountry,
				},
			},
		}
	} else {
		voucherLayout = &epservice.VoucherLayoutFrankingZone
	}

	cartResponse, err := h.service.CheckoutShoppingCartPDF(&epservice.CheckoutShoppingCartPDFRequest{
		UserToken:    authResponse.UserToken,
		PageFormatId: (*epservice.PageFormatId)(&pageFormat),
		Positions: []*epservice.ShoppingCartPDFPosition{{
			ShoppingCartPosition: &epservice.ShoppingCartPosition{
				ProductCode:   (*epservice.ProductCode)(&productCode),
				Address:       addressBinding,
				VoucherLayout: voucherLayout,
			},
			Position: &epservice.VoucherPosition{
				Position: &epservice.Position{
					LabelX: 1,
					LabelY: 1,
				},
				Page: 1,
			},
		}},
		Total:              (*epservice.ShoppingCartPrice)(&total),
		CreateManifest:     (*epservice.Flag)(&createManifest),
		CreateShippingList: &epservice.ShippingList0,
	})

	if err != nil {
		return 0, "", err
	}

	balance := int(*cartResponse.WalletBallance)

	link := string(*cartResponse.Link)
	log.Println("Fetching PDF bytes")
	pdf, err := getStampPDFBytes(link)
	if err != nil {
		return balance, "", err
	}

	//TODO factor out printing
	log.Println("Starting print")
	lprPath, err := exec.LookPath("lpr")

	if err != nil {
		return balance, link, err
	}

	printCmd := exec.Command(lprPath, "-P", PrinterName, "-C", "Internetmarke")
	printCmd.Stdin = bytes.NewBuffer(pdf)

	output, err := printCmd.CombinedOutput()
	if err != nil {
		log.Println("Printing error: ", string(output))
		return balance, link, err
	}

	return int(*cartResponse.WalletBallance), link, nil
}

func processNames(stamp *StampData) (senderCompany, receiverCompany *epservice.CompanyName, senderPerson, receiverPerson *epservice.PersonName) {

	// API demands both given and surname
	splitSenderName := strings.SplitN(stamp.SenderName, " ", 2)
	splitReceiverName := strings.SplitN(stamp.ReceiverName, " ", 2)

	if stamp.SenderCompany != "" {
		senderCompany = &epservice.CompanyName{
			Company: stamp.SenderCompany,
		}
	}

	if stamp.SenderName != "" {
		var s *epservice.PersonName
		if len(splitSenderName) < 2 {
			s = &epservice.PersonName{
				Firstname: " ",
				Lastname:  splitSenderName[0],
			}
		} else {
			s = &epservice.PersonName{
				Firstname: splitSenderName[0],
				Lastname:  splitSenderName[1],
			}
		}
		if senderCompany != nil {
			senderCompany.PersonName = s
		} else {
			senderPerson = s
		}
	}

	if stamp.ReceiverCompany != "" {
		receiverCompany = &epservice.CompanyName{
			Company: stamp.ReceiverCompany,
		}
	}
	if stamp.ReceiverName != "" {
		var s *epservice.PersonName
		if len(splitReceiverName) < 2 {
			s = &epservice.PersonName{
				Firstname: " ",
				Lastname:  splitReceiverName[0],
			}
		} else {
			s = &epservice.PersonName{
				Firstname: splitReceiverName[0],
				Lastname:  splitReceiverName[1],
			}
		}
		if receiverCompany != nil {
			receiverCompany.PersonName = s
		} else {
			receiverPerson = s
		}
	}

	return

}

func checkPageFormat(service *epservice.OneClickForAppPortTypeV3, requestedFormat int) (supportsAddress bool, err error) {

	var pageFormatResponse *epservice.RetrievePageFormatsResponse

	pageFormatResponse, err = service.RetrievePageFormats(&epservice.RetrievePageFormatsRequest{})

	if err != nil {
		return
	}

	//TODO output page formats on user request
	/*for _, v := range pageFormatResponse.PageFormat {
		log.Printf("%d - %s", *v.Id, v.Name)
	}

	var out []byte
	out, err = json.MarshalIndent(pageFormatResponse.PageFormat, "", "\t")

	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write([]byte("\n\n"))
	os.Stdout.Write(out)
	os.Stdout.Write([]byte("\n\n"))
	os.Exit(0)*/

	for _, v := range pageFormatResponse.PageFormat {
		id := int(*v.Id)
		if id == requestedFormat {
			log.Printf("Using page format: %d - %s. Address possible: %v\n", id, v.Name, v.IsAddressPossible)
			return v.IsAddressPossible, nil
		}
	}

	return false, errors.New("Requested page format not found")

}
