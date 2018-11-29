package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/headzoo/surf.v1"

	"github.com/PuerkitoBio/goquery"
)

var (
	_bcpa Bcpa
	_url  string
)

// Bcpa table contains the information for each user
type Bcpa struct {
	Siteaddress         string `json:"siteaddress"`
	Owner               string `json:"owner"`
	MailingAddress      string `json:"mailingAddress"`
	ID                  string `json:"id"`
	Milage              string `json:"milage"`
	Use                 string `json:"use"`
	Legal               string `json:"legal"`
	PropertyAssessments []PropertyAssessmentValue
	ExemptionsTaxable   ExemptionsTaxableValuesbyTaxingAuthority
	SalesHistory        []Sale
	LandCalculations    LandCalculations
	SpecialAssessments  []SpecialAssessment
}

// RecBuildingCard
type RecBuildingCard struct {
	ParcelIDNumber            string `json:"parcelidnumber"`
	UseCode                   string `json:"usecode"`
	NoBedrooms                string `json:"nobedrooms"`
	NoBaths                   string `json:"nobaths"`
	NoUnits                   string `json:"nounits"`
	NoStories                 string `json:"nostories"`
	NoBuildings               string `json:"nobuildings"`
	Foundation                string `json:"foundation"`
	Exterior                  string `json:"exterior"`
	RoofType                  string `json:"rooftype"`
	RoofMaterial              string `json:"roofmaterial"`
	Interior                  string `json:"interior"`
	Floors                    string `json:"floors"`
	Plumbing                  string `json:"plumbing"`
	Electric                  string `json:"electric"`
	Classification            string `json:"classification"`
	CeilingHeights            string `json:"ceilingheights"`
	QualityOfConstruction     string `json:"qualityofconstruction"`
	CurrentConditionStructure string `json:"currentconditionstructure"`
	ConstructionClass         string `json:"constructionclass"`
	Permits                   []Permit
	ExtraFeatures             []ExtraFeature
}

// ExtraFeature
type ExtraFeature struct {
	Feature string `json:"feature"`
}

// Permit
type Permit struct {
	PermitNo   string `json:"permitco"`
	PermitType string `json:"permittype"`
	EstCost    string `json:"estcost"`
	PermitDate string `json:"permitdate"`
	CODate     string `json:"codate"`
}

// LandCalculations
type LandCalculations struct {
	Calculations    []LandCalculation
	AdjBldgSF       string `json:"adjbldgsf"`
	Units           string `json:"units"`
	Card            RecBuildingCard
	SketchURL       string `json:"units"`
	EffActYearBuilt string `json:"effactyearbuilt"`
}

// LandCalculation
type LandCalculation struct {
	Price  string `json:"price"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

// SpecialAssessment
type SpecialAssessment struct {
	Fire  string `json:"fire"`
	Garb  string `json:"garb"`
	Light string `json:"light"`
	Drain string `json:"drain"`
	Impr  string `json:"impr"`
	Safe  string `json:"safe"`
	Storm string `json:"storm"`
	Clean string `json:"clean"`
	Misc  string `json:"misc"`
}

// RecPatriotSketch
type RecPatriotSketch struct {
	Sketch       string `json:"sketch"`
	Building     string `json:"building"`
	URL          string `json:"url"`
	SketchImgURL string `json:"sketchimgurl"`
	Codes        []PatriotSketchCode
	AdjAreaTotal string `json:"adjareatotal"`
}

// PatriotSketchCode
type PatriotSketchCode struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Area        string `json:"area"`
	Factor      string `json:"factor"`
	AdjArea     string `json:"adjarea"`
	Stories     string `json:"stories"`
}

// ExemptionsTaxableValuesbyTaxingAuthority table contains the exemptions
type ExemptionsTaxableValuesbyTaxingAuthority struct {
	County      ExemptionsAndTaxableValue
	SchoolBoard ExemptionsAndTaxableValue
	Municipal   ExemptionsAndTaxableValue
	Independent ExemptionsAndTaxableValue
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
}

// ExemptionsAndTaxableValue table contains the exemption values
type ExemptionsAndTaxableValue struct {
	JustValue    string `json:"justvalue"`
	Portability  string `json:"portability"`
	AssessedSOH  string `json:"assessedsoh"`
	Homestead    string `json:"homestead"`
	AddHomestead string `json:"addhomestead"`
	WidVetDis    string `json:"widvetdis"`
	Senior       string `json:"senior"`
	XemptType    string `json:"xempttype"`
	Taxable      string `json:"taxable"`
}

// PropertyAssessmentValue table contains the house values
type PropertyAssessmentValue struct {
	Year                string    `json:"year"`
	Land                string    `json:"land"`
	BuildingImprovement string    `json:"buildingimprovement"`
	JustMarketValue     string    `json:"justmarketvalue"`
	AssessedSOHValue    string    `json:"assessedsohvalue"`
	Tax                 string    `json:"tax"`
	CreatedAt           time.Time `json:"createdat"`
	UpdatedAt           time.Time `json:"updatedat"`
}

//Sale
type Sale struct {
	Date        string `json:"date"`
	Type        string `json:"type"`
	Price       string `json:"price"`
	BookPageCIN string `json:"bookpagecin"`
}

// ParseRecord table contains the information for each user
func ParseRecord(s *goquery.Selection) PropertyAssessmentValue {
	p := PropertyAssessmentValue{}

	// Loop through each cell
	s.Find("td").Each(func(int int, s *goquery.Selection) {

		switch int {
		case 0:
			p.Year = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 1:
			p.Land = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 2:
			p.BuildingImprovement = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 3:
			p.JustMarketValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 4:
			p.AssessedSOHValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 5:
			p.Tax = strings.TrimSpace(s.Find("span").First().Contents().Text())
		}
	})

	return p
}

func main() {
	// Create a new browser and open reddit.
	bow := surf.NewBrowser()
	err := bow.Open("http://www.bcpa.net/RecAddr.asp")
	if err != nil {
		panic(err)
	}

	// Log in to the site.
	fm, _ := bow.Form("[name='homeind']")
	fm.Input("Situs_Street_Number", "515")
	fm.SelectByOptionValue("Situs_Street_Direction", "SW")
	fm.Input("Situs_Street_Name", "18")
	fm.SelectByOptionValue("Situs_Street_Type", "AVE")
	fm.Input("Situs_Street_Post_Dir", "")
	fm.Input("Situs_Unit_Number", "15")
	fm.SelectByOptionValue("Situs_City", "FL")

	if fm.Submit() != nil {
		panic(err)
	}

	//fmt.Println(bow.Body())

	//fmt.Println(bow.Url())

	// Load the HTML document from the URL
	doc, err := goquery.NewDocument(bow.Url().String())
	if err != nil {
		log.Fatal(err)
	}

	//Load the BCPA parent node from the HTML receieved from URL
	_bcpa, err = LoadBcpaFromDoc(doc.Html())

	if err != nil {
		log.Fatal(err)
	}

	//doc.Html().string
	//Load the class level BCPA object with with assessments
	err = LoadAppendPropertyAssessments(doc.Html())

	if err != nil {
		log.Fatal(err)
	}

	//load exemptions
	err = LoadAppendExemptionsTaxable(doc.Html())

	if err != nil {
		log.Fatal(err)
	}

	//Load Sales History
	err = LoadSalesHistory(doc.Html())

	if err != nil {
		log.Fatal(err)
	}

	//log.Println(len(_bcpa.PropertyAssessments))
	//log.Println(_bcpa)

	file, err := os.Create("C:\\gowork\\testFiles\\result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	file.Write([]byte(marshalBcpa(_bcpa)))

	//fmt.Fprintf(file, marshalBcpa(_bcpa))

}

// ParseSalesRecord Parse Sales hostory table
func ParseSalesRecord(s *goquery.Selection) Sale {
	sale := Sale{}

	s.Find("td").Each(func(int int, s *goquery.Selection) {

		switch int {
		case 0:
			sale.Date = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 1:
			sale.Type = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 2:
			sale.Price = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 3:
			sale.BookPageCIN = strings.TrimSpace(s.Find("span").First().Contents().Text())
		}

	})

	return sale
}

// LoadSalesHistory Load up the sales history table in objects and append to BCPA parent
func LoadSalesHistory(html string, e error) error {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}
	//Preload the object
	//sale := Sale{}

	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(10) > tbody > tr > td:nth-child(1) > table:nth-child(1) > tbody > tr").Each(func(i int, s *goquery.Selection) {

		if i > 1 {

			sale := ParseSalesRecord(s)
			_bcpa.SalesHistory = append(_bcpa.SalesHistory, sale)
		}
	})

	return err
}

// ParseExemptionsTaxableRecord ...
func ParseExemptionsTaxableRecord(s *goquery.Selection, i int, eta ExemptionsTaxableValuesbyTaxingAuthority) ExemptionsTaxableValuesbyTaxingAuthority {

	// Loop through each cell
	s.Find("td").Each(func(int int, s *goquery.Selection) {
		//fmt.Println(i)
		//fmt.Println(strings.TrimSpace(s.Find("span").First().Contents().Text()))

		switch i {
		case 2:
			switch int {
			case 1:
				//p.Land = strings.TrimSpace(s.Find("span").First().Contents().Text())
				eta.County.JustValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.JustValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.JustValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.JustValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 3:
			switch int {
			case 1:
				eta.County.Portability = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.Portability = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.Portability = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.Portability = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 4:
			switch int {
			case 1:
				eta.County.AssessedSOH = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.AssessedSOH = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.AssessedSOH = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.AssessedSOH = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 5:
			switch int {
			case 1:
				eta.County.Homestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.Homestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.Homestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.Homestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 6:
			switch int {
			case 1:
				eta.County.AddHomestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.AddHomestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.AddHomestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.AddHomestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 7:
			switch int {
			case 1:
				eta.County.WidVetDis = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.WidVetDis = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.WidVetDis = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.WidVetDis = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 8:
			switch int {
			case 1:
				eta.County.Senior = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.Senior = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.Senior = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.Senior = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 9:
			switch int {
			case 1:
				eta.County.XemptType = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.XemptType = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.XemptType = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.XemptType = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 10:
			switch int {
			case 1:
				eta.County.Taxable = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.Taxable = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.Taxable = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.Taxable = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		}

	})

	return eta
}

// LoadAppendExemptionsTaxable ...
func LoadAppendExemptionsTaxable(html string, e error) error {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}
	//Preload the object
	eta := ExemptionsTaxableValuesbyTaxingAuthority{}
	eta.CreatedAt = time.Now()
	eta.County = ExemptionsAndTaxableValue{}
	eta.SchoolBoard = ExemptionsAndTaxableValue{}
	eta.Municipal = ExemptionsAndTaxableValue{}
	eta.Independent = ExemptionsAndTaxableValue{}

	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(8) > tbody > tr").Each(func(i int, s *goquery.Selection) {

		if i > 1 {

			eta = ParseExemptionsTaxableRecord(s, i, eta)
			//fmt.Println(eta)

		}
	})

	_bcpa.ExemptionsTaxable = eta

	return err
}

//LoadAppendPropertyAssessments used to load and append Assessments to the BCPA parent node
func LoadAppendPropertyAssessments(html string, e error) error {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(6) > tbody > tr").Each(func(i int, s *goquery.Selection) {

		if i > 1 {
			pa := ParseRecord(s)
			pa.CreatedAt = time.Now()
			_bcpa.PropertyAssessments = append(_bcpa.PropertyAssessments, pa)
		}
	})

	return err
}

//LoadBcpaFromDoc used to load Bcpa data from HTML
func LoadBcpaFromDoc(html string, e error) (Bcpa, error) {

	var bcpa Bcpa
	var siteAddress, owner, mailingAddress, id, mileage, use, legal string

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	// use selector found with the browser inspector
	siteAddress = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(1) > td:nth-child(2) > span > a > b").Contents().Text()

	//clean up the carriage return
	re := regexp.MustCompile(`\r?\n`)
	siteAddress = re.ReplaceAllString(siteAddress, " ")
	//siteAddress = strings.Replace(siteAddress, " 			  ", " ", 1)

	releadclosewhtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reinsidewhtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := releadclosewhtsp.ReplaceAllString(siteAddress, "")
	siteAddress = reinsidewhtsp.ReplaceAllString(final, " ")

	//Set the Object
	bcpa.Siteaddress = strings.TrimSpace(siteAddress)

	owner = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(2) > td:nth-child(2) > span").Contents().Text()
	//Set the Object
	bcpa.Owner = strings.TrimSpace(owner)

	mailingAddress = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(3) > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.MailingAddress = strings.TrimSpace(mailingAddress)

	id = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(1) > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.ID = strings.TrimSpace(strings.Replace(StripSpaces(id), " ", "", -1))

	mileage = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(2) > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.Milage = strings.TrimSpace(mileage)

	use = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(3) > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.Use = strings.TrimSpace(StripSpaces(use))

	legal = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(4) > tbody > tr > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.Legal = strings.TrimSpace(legal)

	return bcpa, err
}

//StripSpaces remove leading and trailing and extra gapped spaces
func StripSpaces(o string) string {

	releadclosewhtsp2 := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)

	reinsidewhtsp2 := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	final := releadclosewhtsp2.ReplaceAllString(o, "")

	return reinsidewhtsp2.ReplaceAllString(final, " ")
}

func marshalBcpa(bcpa Bcpa) string {
	//user := &User{name:"Frank"}
	b, err := json.Marshal(bcpa)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "0"
	}
	fmt.Println(string(b))

	return string(b)
}
