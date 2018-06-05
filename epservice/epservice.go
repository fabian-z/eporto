// TODO test remaining unused methods
// TODO remove XML names for nested elements

// Base package generated by github.com/hooklift/gowsdl
package epservice

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type VoucherLayout string

var (
	VoucherLayoutAddressZone  VoucherLayout = "AddressZone"
	VoucherLayoutFrankingZone VoucherLayout = "FrankingZone"
)

// Ja/Nein-Entscheidung.
type Flag bool

// Versandliste Flag mögliche Werte:	0 = keine Versandliste; 1 = Versandliste ohne Adressen;	2 = Versandliste mit Adressen
type ShippingList string

var (
	ShippingList0 ShippingList = "0"
	ShippingList1 ShippingList = "1"
	ShippingList2 ShippingList = "2"
)

// ID des Druckformates
type PageFormatId int32

// Das aktuelle Guthaben der Portokasse des Kunden in der  PC-Frankierung, Betrag in Cent
type WalletBalance int32

// Definiert eine feste Liste von ErrorCodes die fachlichen Fehlern bei der Prüfung des Warenkorbs auf den Systemen der PC-Frankierung entsprechen.
type AuthenticateUserErrorCodes string

const (
	AuthenticateUserErrorCodesUnkownUser  AuthenticateUserErrorCodes = "unkownUser"
	AuthenticateUserErrorCodesInvalidUser AuthenticateUserErrorCodes = "invalidUser"
)

// Definiert eine feste Liste von ErrorCodes die fachlichen Fehlern bei der Prüfung des Warenkorbs auf den Systemen der PC-Frankierung entsprechen.
type ShoppingCartValidationErrorCodes string

const (
	ShoppingCartValidationErrorCodesInvalidUser ShoppingCartValidationErrorCodes = "invalidUser"

	ShoppingCartValidationErrorCodesInvalidPplId ShoppingCartValidationErrorCodes = "invalidPplId"

	ShoppingCartValidationErrorCodesInvalidProductcode ShoppingCartValidationErrorCodes = "invalidProductcode"

	ShoppingCartValidationErrorCodesInvalidTotalAmount ShoppingCartValidationErrorCodes = "invalidTotalAmount"

	ShoppingCartValidationErrorCodesWalletBalanceNotEnough ShoppingCartValidationErrorCodes = "walletBalanceNotEnough"

	ShoppingCartValidationErrorCodesWalletNotAvailable ShoppingCartValidationErrorCodes = "walletNotAvailable"

	ShoppingCartValidationErrorCodesInvalidMotive ShoppingCartValidationErrorCodes = "invalidMotive"

	ShoppingCartValidationErrorCodesInvalidPageFormat ShoppingCartValidationErrorCodes = "invalidPageFormat"

	ShoppingCartValidationErrorCodesProductExpired ShoppingCartValidationErrorCodes = "productExpired"

	ShoppingCartValidationErrorCodesInvalidShopOrderId ShoppingCartValidationErrorCodes = "invalidShopOrderId"

	ShoppingCartValidationErrorCodesInvalidOrderPositionCount ShoppingCartValidationErrorCodes = "invalidOrderPositionCount"
)

// ID eines Produktes
type ProductCode int32

// Preis eines Produktes in Eurocent
type ProductPrice int32

// PPL (Product Price List)
type PPL int32

// ID eines Bildes aus der Motivgallery
type ImageID int32

// Wert des Warenkorbs in Eurocent
type ShoppingCartPrice int32

// URL zum Druck der Warenkorbpositionen
type Link string

// Benutzerkennung
type UserToken string

type ShopOrderId string

// Definiert eine feste Liste von ErrorCodes die fachlichen Fehlern beim Einstellen und Abfragen von Orders des Frankierbackends entsprechen.
type RetrieveOrderErrorCodes string

const (
	RetrieveOrderErrorCodesUnknownShopOrderId RetrieveOrderErrorCodes = "unknownShopOrderId"
)

// Ausrichtung
type Orientation string

const (
	OrientationLANDSCAPE Orientation = "LANDSCAPE"

	OrientationPORTRAIT Orientation = "PORTRAIT"
)

// ENUM für alle möglichen Typen von Druckmedien.
type PageType string

const (
	PageTypeREGULARPAGE PageType = "REGULARPAGE"

	PageTypeENVELOPE PageType = "ENVELOPE"

	PageTypeLABELPRINTER PageType = "LABELPRINTER"

	PageTypeLABELPAGE PageType = "LABELPAGE"
)

type CreateShopOrderIdRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:CreateShopOrderIdRequest"`

	UserToken *UserToken `xml:"userToken,omitempty"`
}

type CreateShopOrderIdResponse struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 CreateShopOrderIdResponse"`

	ShopOrderId *ShopOrderId `xml:"shopOrderId,omitempty"`
}

type AuthenticateUserRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:AuthenticateUserRequest"`

	// Benutzername (Email) des IM-Benutzers.
	Username string `xml:"ns0:username,omitempty"`

	// Passwort, das bei der Registrierung des IM-Benutzers vergeben wurde.
	Password string `xml:"ns0:password,omitempty"`
}

type AuthenticateUserResponse struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 AuthenticateUserResponse"`

	// Token zur Benutzer Indentifizierung
	//
	UserToken *UserToken `xml:"userToken,omitempty"`

	// Aktuelles Guthaben der Portokasse in Eurocent
	//
	WalletBalance *WalletBalance `xml:"walletBalance,omitempty"`

	// TRUE, falls der Benutzer den aktuellen
	// AGB's zugestimmt hat, ansonst FALSE
	//
	ShowTermsAndConditions bool `xml:"showTermsAndConditions,omitempty"`

	InfoMessage string `xml:"infoMessage,omitempty"`
}

type RetrievePreviewVoucherPDFRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:RetrievePreviewVoucherPDFRequest"`

	// ID die das Produkt innerhalb der PPL identifiziert
	ProductCode *ProductCode `xml:"productCode,omitempty"`

	// ID des Motives das am Vorschaubild dargestellt werden soll
	ImageID *ImageID `xml:"imageID,omitempty"`

	// Layout des Frankiervermerks zB Adressgebunden, Adressungebunden
	VoucherLayout *VoucherLayout `xml:"voucherLayout,omitempty"`

	// ID des Druckformates
	PageFormatId *PageFormatId `xml:"pageFormatId,omitempty"`
}

type RetrievePreviewVoucherPNGRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:RetrievePreviewVoucherPNGRequest"`

	// ID die das Produkt innerhalb der PPL identifiziert
	ProductCode *ProductCode `xml:"productCode,omitempty"`

	// ID des Motives das am Vorschaubild dargestellt werden soll
	ImageID *ImageID `xml:"imageID,omitempty"`

	// Layout des Frankiervermerks zB Adressgebunden, Adressungebunden
	VoucherLayout *VoucherLayout `xml:"voucherLayout,omitempty"`
}

type RetrievePreviewVoucherResponse struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 RetrievePreviewVoucherResponse"`

	// Link auf eine PNG/PDF Datei
	Link *Link `xml:"link,omitempty"`
}

type MotiveLink struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 MotiveLink"`

	// URL zum Motiv
	Link *Link `xml:"link,omitempty"`

	// URL zum Motiv-Vorschaubild (Thumbnail)
	LinkThumbnail *Link `xml:"linkThumbnail,omitempty"`
}

type RetrievePrivateGalleryRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:RetrievePrivateGalleryRequest"`

	// ID die den Benutzer eindeutig identifiziert
	UserToken *UserToken `xml:"userToken,omitempty"`
}

type RetrievePrivateGalleryResponse struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 RetrievePrivateGalleryResponse"`

	ImageLink []*MotiveLink `xml:"imageLink,omitempty"`
}

type ShoppingCartResponse struct {
	//XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 CheckoutShoppingCartResponse"`

	// Ein Link auf die Datei (ZIP Archiv oder PDF) mit den Internetmarken.
	Link *Link `xml:"link,omitempty"`

	// Ein Link auf die Datei (PDF) mit dem Einlieferbeleg / der Versandliste.
	ManifestLink *Link `xml:"manifestLink,omitempty"`

	// Portokassenguthaben nach Kauf
	WalletBallance *WalletBalance `xml:"walletBallance,omitempty"`

	ShoppingCart *ShoppingCart `xml:"shoppingCart,omitempty"`
}

type ShoppingCart struct {
	//XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 ShoppingCart"`

	ShopOrderId *ShopOrderId `xml:"shopOrderId,omitempty"`

	VoucherList *VoucherList `xml:"voucherList,omitempty"`
}

type VoucherList struct {
	//XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 VoucherList"`

	Voucher []*VoucherType `xml:"voucher,omitempty"`
}

type CheckoutShoppingCartPNGRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:CheckoutShoppingCartPNGRequest"`

	// Eindeutige Benutzerkennung
	UserToken *UserToken `xml:"userToken,omitempty"`

	ShopOrderId *ShopOrderId `xml:"shopOrderId,omitempty"`

	// PPL auf die sich die Produkte im Warenkorb beziehen
	Ppl *PPL `xml:"ppl,omitempty"`

	// Warenkorbpositionen
	Positions []*ShoppingCartPosition `xml:"positions,omitempty"`

	// Gesamtwert des Warenkorbs
	Total *ShoppingCartPrice `xml:"total,omitempty"`

	// Flag, das anzeigt, ob ein Einlieferungsbeleg erstellt werden soll.
	CreateManifest *Flag `xml:"createManifest,omitempty"`

	// Enum, welche festlegt, ob eine Versandliste erstellt werden soll und wenn ja, ob mit oder ohne Adressen.
	CreateShippingList *ShippingList `xml:"createShippingList,omitempty"`
}

type CheckoutShoppingCartPDFRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:CheckoutShoppingCartPDFRequest"`

	// Eindeutige Benutzerkennung
	UserToken *UserToken `xml:"ns0:userToken,omitempty"`

	ShopOrderId *ShopOrderId `xml:"ns0:shopOrderId,omitempty"`

	// ID des Druckformates auf das der Ausdruck erfolgt.
	PageFormatId *PageFormatId `xml:"ns0:pageFormatId,omitempty"`

	// PPL auf die sich die Produkte im Warenkorb beziehen
	Ppl *PPL `xml:"ns0:ppl,omitempty"`

	// Warenkorbpositionen
	Positions []*ShoppingCartPDFPosition `xml:"ns0:positions,omitempty"`

	// Gesamtwert des Warenkorbs
	Total *ShoppingCartPrice `xml:"ns0:total,omitempty"`

	// Flag, das anzeigt, ob ein Einlieferungsbeleg erstellt werden soll.
	CreateManifest *Flag `xml:"ns0:createManifest,omitempty"`

	// Enum, welche festlegt, ob eine Versandliste erstellt werden soll und wenn ja, ob mit oder ohne Adressen.
	CreateShippingList *ShippingList `xml:"ns0:createShippingList,omitempty"`
}

type ShoppingCartValidationErrorInfo struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 ShoppingCartValidationErrorInfo"`

	Id *ShoppingCartValidationErrorCodes `xml:"id,omitempty"`

	Message string `xml:"message,omitempty"`
}

type AuthenticateUserException struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 AuthenticateUserException"`

	Id *AuthenticateUserErrorCodes `xml:"id,omitempty"`

	Message string `xml:"message,omitempty"`
}

type IdentifyException struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 IdentifyException"`

	Message string `xml:"message,omitempty"`
}

type InvalidProductException struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 InvalidProductException"`

	Message string `xml:"message,omitempty"`
}

type InvalidPageFormatException struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 InvalidPageFormatException"`

	Message string `xml:"message,omitempty"`
}

type InvalidMotiveException struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 InvalidMotiveException"`

	Message string `xml:"message,omitempty"`
}

type ShoppingCartValidationException struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 ShoppingCartValidationException"`

	Message string `xml:"message,omitempty"`

	Errors []*ShoppingCartValidationErrorInfo `xml:"errors,omitempty"`
}

type RetrievePublicGalleryRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:RetrievePublicGalleryRequest"`
}

type RetrievePublicGalleryResponse struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 RetrievePublicGalleryResponse"`

	Items []*GalleryItem `xml:"items,omitempty"`
}

type ImageItem struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 ImageItem"`

	// Name des Motives
	ImageID *ImageID `xml:"imageID,omitempty"`

	// Beschreibung des Motives
	ImageDescription string `xml:"imageDescription,omitempty"`

	// Kurzer Slogan zum Bild
	ImageSlogan string `xml:"imageSlogan,omitempty"`

	// URL zum Bild
	Links *MotiveLink `xml:"links,omitempty"`
}

type GalleryItem struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 GalleryItem"`

	// Motivkategorie
	Category string `xml:"category,omitempty"`

	// Motivkategorie Beschreibung
	CategoryDescription string `xml:"categoryDescription,omitempty"`

	// Motivkategorie ID
	CategoryId int32 `xml:"categoryId,omitempty"`

	// Liste aller Bilder innerhalb der Bilderkategorie
	Images []*ImageItem `xml:"images,omitempty"`
}

type Name struct {
	//XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:Name"`

	PersonName *PersonName `xml:"ns0:personName,omitempty"`

	CompanyName *CompanyName `xml:"ns0:companyName,omitempty"`
}

type PersonName struct {
	//XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:PersonName"`

	// Anrede
	Salutation string `xml:"ns0:salutation,omitempty"`

	// Titel
	Title string `xml:"ns0:title,omitempty"`

	// Vorname
	Firstname string `xml:"ns0:firstname,omitempty"`

	// Nachname
	Lastname string `xml:"ns0:lastname,omitempty"`
}

type CompanyName struct {
	//XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:CompanyName"`

	// Firmenname
	Company string `xml:"ns0:company,omitempty"`

	PersonName *PersonName `xml:"ns0:personName,omitempty"`
}

type Address struct {
	//XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:Address"`

	// Adresszusatz, z.B. Im Hinterhof 2. Tür rechts
	Additional string `xml:"ns0:additional,omitempty"`

	// Straße
	Street string `xml:"ns0:street,omitempty"`

	// Hausnummer
	HouseNo string `xml:"ns0:houseNo,omitempty"`

	// Postleitzahl
	Zip string `xml:"ns0:zip,omitempty"`

	// Name des Ortes
	City string `xml:"ns0:city,omitempty"`

	// 3-stelliger ISO-Ländercode
	Country string `xml:"ns0:country,omitempty"`
}

type NamedAddress struct {
	//XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:NamedAddress"`

	Name *Name `xml:"ns0:name,omitempty"`

	Address *Address `xml:"ns0:address,omitempty"`
}

type ShoppingCartPosition struct {
	//XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 ShoppingCartPosition"`

	// ID des Produktes
	ProductCode *ProductCode `xml:"ns0:productCode,omitempty"`

	// Falls ein Motiv vorhanden, dann ID übergeben
	ImageID *ImageID `xml:"ns0:imageID,omitempty"`

	// Bei adressgebundenen Produkten muss der Sender und Empfänger angegeben werden
	Address *AddressBinding `xml:"ns0:address,omitempty"`

	// Zusatzinfo zur Bestellposition
	AdditionalInfo string `xml:"ns0:additionalInfo,omitempty"`

	// Bestimmt das Layout des Frankiervermerks.
	//
	VoucherLayout *VoucherLayout `xml:"ns0:voucherLayout,omitempty"`
}

type AddressBinding struct {
	//XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:AddressBinding"`

	Sender *NamedAddress `xml:"ns0:sender,omitempty"`

	Receiver *NamedAddress `xml:"ns0:receiver,omitempty"`
}

type RetrieveOrderException struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 RetrieveOrderException"`

	Message string `xml:"message,omitempty"`

	Errors []*RetrieveOrderErrorCodes `xml:"errors,omitempty"`
}

type RetrieveOrderRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:RetrieveOrderRequest"`

	// Eindeutige Benutzerkennung
	UserToken *UserToken `xml:"ns0:userToken,omitempty"`

	// Shop Order-ID.
	ShopOrderId *ShopOrderId `xml:"ns0:shopOrderId,omitempty"`
}

type RetrieveOrderResponse struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 RetrieveOrderResponse"`

	// Link zum Druck des Warenkorbs
	Link *Link `xml:"link,omitempty"`

	// Link zum Druck des Einlieferungsbelegs/Versandliste
	ManifestLink *Link `xml:"manifestLink,omitempty"`

	// Detaillierte Beschreibung des erzeugten Warenkorbs
	ShoppingCart *ShoppingCart `xml:"shoppingCart,omitempty"`
}

type VoucherPosition struct {
	//XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 VoucherPosition"`

	*Position

	// Auf welcher Seite befindet sich die Marke
	//
	Page int32 `xml:"ns0:page,omitempty"`
}

type Position struct {
	//XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 Position"`

	LabelX int32 `xml:"ns0:labelX,omitempty"`

	LabelY int32 `xml:"ns0:labelY,omitempty"`
}

type ShoppingCartPDFPosition struct {
	//XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 ShoppingCartPDFPosition"`

	*ShoppingCartPosition

	Position *VoucherPosition `xml:"ns0:position,omitempty"`
}

type VoucherType struct {
	//XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 VoucherType"`

	VoucherId string `xml:"voucherId,omitempty"`

	TrackId string `xml:"trackId,omitempty"`
}

type RetrievePageFormatsRequest struct {
	XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 ns0:RetrievePageFormatsRequest"`
}

type RetrievePageFormatsResponse struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 RetrievePageFormatsResponse"`

	PageFormat []*PageFormat `xml:"pageFormat,omitempty"`
}

type PageFormat struct {
	XMLName xml.Name `xml:"http://oneclickforapp.dpag.de/V3 pageFormat"`

	Id *PageFormatId `xml:"id,omitempty"`

	// ja, wenn mit dem Druckformat Adressen auf die Frankiervermerke gedruckt werden können.
	IsAddressPossible bool `xml:"isAddressPossible,omitempty"`

	// ja, wenn mit dem Druckformat Motive auf die Frankiervermerke gedruckt werden können.
	IsImagePossible bool `xml:"isImagePossible,omitempty"`

	// Der Name des Druckformats, z.B. DIN A4 Normalpapier oder Brief C5 162 x 229
	Name string `xml:"name,omitempty"`

	// Beschreibung des Druckformats.
	Description string `xml:"description,omitempty"`

	// Spezifikation des Druckmediums.
	PageType *PageType `xml:"pageType,omitempty"`

	PageLayout struct {

		// Größe des Druckformates in Millimeter
		//
		Size *Dimension `xml:"size,omitempty"`

		// Seitenausrichtung
		Orientation *Orientation `xml:"orientation,omitempty"`

		// Abstand zwischen den Etiketten in Millimeter
		//
		LabelSpacing *Dimension `xml:"labelSpacing,omitempty"`

		// Anzahl der Labelpositionen in X und Y Richtung
		LabelCount *Position `xml:"labelCount,omitempty"`

		// Innerer Randabstand des Druckformates in Millimeter
		Margin *BorderDimension `xml:"margin,omitempty"`
	} `xml:"pageLayout,omitempty"`
}

type BorderDimension struct {
	//XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 BorderDimension"`

	Top float64 `xml:"top,omitempty"`

	Bottom float64 `xml:"bottom,omitempty"`

	Left float64 `xml:"left,omitempty"`

	Right float64 `xml:"right,omitempty"`
}

type Dimension struct {
	//XMLName xml.Name `xml:"ns0:http://oneclickforapp.dpag.de/V3 Dimension"`

	X float64 `xml:"x,omitempty"`

	Y float64 `xml:"y,omitempty"`
}

type OneClickForAppPortTypeV3 struct {
	client *SOAPClient
}

func NewOneClickForAppPortTypeV3(url string, tls bool, auth *AuthenticationData) (*OneClickForAppPortTypeV3, error) {
	if url == "" {
		url = ""
	}

	if auth == nil {
		return nil, errors.New("No authentiction data")
	}

	client := NewSOAPClient(url, auth)

	return &OneClickForAppPortTypeV3{
		client: client,
	}, nil
}

func NewOneClickForAppPortTypeV3WithTLSConfig(url string, tlsCfg *tls.Config, auth *AuthenticationData) (*OneClickForAppPortTypeV3, error) {
	if url == "" {
		url = ""
	}

	if auth == nil {
		return nil, errors.New("No authentiction data")
	}

	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &OneClickForAppPortTypeV3{
		client: client,
	}, nil
}

func (service *OneClickForAppPortTypeV3) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *OneClickForAppPortTypeV3) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

func (service *OneClickForAppPortTypeV3) RetrievePublicGallery(request *RetrievePublicGalleryRequest) (*RetrievePublicGalleryResponse, error) {
	response := new(RetrievePublicGalleryResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - IdentifyException
//   - ShoppingCartValidationException

func (service *OneClickForAppPortTypeV3) CheckoutShoppingCartPDF(request *CheckoutShoppingCartPDFRequest) (*ShoppingCartResponse, error) {
	response := new(ShoppingCartResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - IdentifyException
//   - ShoppingCartValidationException

func (service *OneClickForAppPortTypeV3) CheckoutShoppingCartPNG(request *CheckoutShoppingCartPNGRequest) (*ShoppingCartResponse, error) {
	response := new(ShoppingCartResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - AuthenticateUserException

func (service *OneClickForAppPortTypeV3) AuthenticateUser(request *AuthenticateUserRequest) (*AuthenticateUserResponse, error) {
	response := new(AuthenticateUserResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidProductException
//   - InvalidMotiveException
//   - InvalidPageFormatException

func (service *OneClickForAppPortTypeV3) RetrievePreviewVoucherPDF(request *RetrievePreviewVoucherPDFRequest) (*RetrievePreviewVoucherResponse, error) {
	response := new(RetrievePreviewVoucherResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - InvalidProductException
//   - InvalidMotiveException

func (service *OneClickForAppPortTypeV3) RetrievePreviewVoucherPNG(request *RetrievePreviewVoucherPNGRequest) (*RetrievePreviewVoucherResponse, error) {
	response := new(RetrievePreviewVoucherResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - IdentifyException

func (service *OneClickForAppPortTypeV3) RetrievePrivateGallery(request *RetrievePrivateGalleryRequest) (*RetrievePrivateGalleryResponse, error) {
	response := new(RetrievePrivateGalleryResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - IdentifyException
//   - RetrieveOrderException

func (service *OneClickForAppPortTypeV3) RetrieveOrder(request *RetrieveOrderRequest) (*RetrieveOrderResponse, error) {
	response := new(RetrieveOrderResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - IdentifyException

func (service *OneClickForAppPortTypeV3) CreateShopOrderId(request *CreateShopOrderIdRequest) (*CreateShopOrderIdResponse, error) {
	response := new(CreateShopOrderIdResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *OneClickForAppPortTypeV3) RetrievePageFormats(request *RetrievePageFormatsRequest) (*RetrievePageFormatsResponse, error) {
	response := new(RetrievePageFormatsResponse)
	err := service.client.Call(" ", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *SOAPHeader
	Body    SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name `xml:"Header"`

	PartnerId        string `xml:"v3:http://oneclickforapp.dpag.de v3:PARTNER_ID"`
	RequestTimestamp string `xml:"v3:http://oneclickforapp.dpag.de v3:REQUEST_TIMESTAMP"`
	KeyPhase         string `xml:"v3:http://oneclickforapp.dpag.de v3:KEY_PHASE"`
	PartnerSignature string `xml:"v3:http://oneclickforapp.dpag.de v3:PARTNER_SIGNATURE"`
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

type SOAPClient struct {
	url     string
	tlsCfg  *tls.Config
	auth    *AuthenticationData
	headers []interface{}
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, auth *AuthenticationData) *SOAPClient {
	return NewSOAPClientWithTLSConfig(url, nil, auth)
}

func NewSOAPClientWithTLSConfig(url string, tlsCfg *tls.Config, auth *AuthenticationData) *SOAPClient {
	return &SOAPClient{
		url:    url,
		tlsCfg: tlsCfg,
		auth:   auth,
	}
}

func (s *SOAPClient) AddHeader(header interface{}) {
	s.headers = append(s.headers, header)
}

// Namespace mapping, see https://github.com/golang/go/issues/13400
var nsReplacer = strings.NewReplacer(`xmlns="v3:http://oneclickforapp.dpag.de"`, `xmlns:v3="http://oneclickforapp.dpag.de"`,
	`xmlns="ns0:http://oneclickforapp.dpag.de/V3"`, `xmlns:ns0="http://oneclickforapp.dpag.de/V3"`)

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{}

	ts := get1C4ATimestamp()
	sig := calculate1C4ASignature(s.auth, ts)

	envelope.Header = &SOAPHeader{
		PartnerId:        s.auth.PartnerId,
		RequestTimestamp: ts,
		KeyPhase:         s.auth.KeyPhase,
		PartnerSignature: sig,
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	requestBody := bytes.NewBufferString("<?xml version='1.0' encoding='utf-8'?>")
	nsReplacer.WriteString(requestBody, buffer.String())

	log.Println("Got here, writing:\n", requestBody.String())

	req, err := http.NewRequest("POST", s.url, requestBody)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", soapAction)

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: s.tlsCfg,
		Dial:            dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		log.Println("empty response")
		return nil
	}
	log.Println("Returned")

	log.Println(string(rawbody))
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}

// CUSTOM

type AuthenticationData struct {
	PartnerId string
	KeyPhase  string
	Key       string
}

func calculate1C4ASignature(a *AuthenticationData, requestTimestamp string) string {
	hash := md5.New()
	_, err := fmt.Fprintf(hash, "%s::%s::%s::%s", a.PartnerId, requestTimestamp, a.KeyPhase, a.Key)

	if err != nil {
		log.Fatal(hash)
	}

	return hex.EncodeToString(hash.Sum([]byte{}))[:8]
}

func get1C4ATimestamp() string {
	return time.Now().Format("02012006-150405")
}
