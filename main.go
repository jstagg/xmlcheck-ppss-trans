package main

/**
Use case is that I had a bunch of existing XML files
of "unknown" correctness (as in: some were only half complete)
that also needed a new element/value added that fortunately
had to match an existing value.

So, "invalid" XML files are suffixed as .BAD to keep them
from being picked up and to make them easy to spot for a
human who is looking for them.

Then the value of Day is copied into WK_DAY to save a bunch
of manual ETL re-work.

Finally, the big, ugly struct for Record is making most of you
SMH right now. I didn't need the full struct, but I wanted to
get the practice of being complete about it. And, yes, this is
how it looks in the XSD (caps and underscores which the linter
absolutely LOVES). So IIWII (It Is What It Is).

And it bears saying: the app loading this cannot tolerate
ANY variance from its XSD. It ain't complicated, but it also
ain't flexible.

**/

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

// Records struct for the transaction set
type Records struct {
	XMLName xml.Name `xml:"TRANSACTION_Set"`
	Record  []Record `xml:"TRANSACTION_Record"`
}

// Record struct for the transaction records
type Record struct {
	XMLName          xml.Name `xml:"TRANSACTION_Record"`
	TRANSACTION_ID   string   `xml:"TRANSACTION_ID"`
	Day              string   `xml:"Day,omitempty"`
	DAY              string   `xml:"DAY"`
	WK_DAY           string   `xml:"WK_DAY"`
	DEALER_ID        string   `xml:"DEALER_ID"`
	CUST_ID          string   `xml:"CUST_ID"`
	PRODUCT_ID       string   `xml:"PRODUCT_ID"`
	CONTRACT         string   `xml:"CONTRACT"`
	SALES_PERSON_ID  string   `xml:"SALES_PERSON_ID"`
	ORDER_TYPE       string   `xml:"ORDER_TYPE"`
	ORDER_ID         string   `xml:"ORDER_ID"`
	ORDER_LINE_ID    string   `xml:"ORDER_LINE_ID"`
	SOLD_UOM         string   `xml:"SOLD_UOM"`
	INVOICE_NUMBER   string   `xml:"INVOICE_NUMBER"`
	HOW_ORDER_PLACED string   `xml:"HOW_ORDER_PLACED"`
	DEALER_FIELD_1   string   `xml:"DEALER_FIELD_1"`
	DEALER_FIELD_2   string   `xml:"DEALER_FIELD_2"`
	QUANTITY         string   `xml:"QUANTITY"`
	INVOICE_SALES    string   `xml:"INVOICE_SALES"`
	UNIT_PRICE       string   `xml:"UNIT_PRICE"`
	COGS             string   `xml:"COGS"`
	FREIGHT          string   `xml:"FREIGHT"`
	REBATES          string   `xml:"REBATES"`
	OTHER_CHARGES    string   `xml:"OTHER_CHARGES"`
	UNIT_COST        string   `xml:"UNIT_COST"`
	LIST_PRICE       string   `xml:"LIST_PRICE"`
	DISCOUNT         string   `xml:"DISCOUNT"`
	SALES_OBJ_ID     string   `xml:"SALES_OBJ_ID"`
	Extraction_Time  string   `xml:"Extraction_Time"`
}

func main() {

	// Wow, this makes setting command line flags easy!
	fptr := flag.String("fp", "TRANSACTION-test.xml", "file path to read")
	flag.Parse()

	// Perl variable habits die hard
	fileOrig := *fptr
	xmlFile, err := os.Open(*fptr)
	badSuffix := ".BAD"
	fileBad := *fptr + badSuffix

	// Perl regex addictions die hard
	// But I don't want to reprocess an existing BAD file
	matched, err := regexp.MatchString(`.*\.BAD`, *fptr)
	if err != nil {
		log.Fatal(err)
	}
	if matched == true {
		log.Fatal("Already a BAD file? ", matched)
	}
	// I mean like REALLY hard
	// But I have to bail on these if that's in the filename
	matched, err2 := regexp.MatchString(`TRANSACTION`, *fptr)
	if err2 != nil {
		log.Fatal(err2)
	}
	if matched == false {
		log.Fatal("Is it a transaction file? ", matched)
	}

	// And if I can't open it for some other reason, assume it's bad
	if err != nil {
		fmt.Println(err)
		err := os.Rename(fileOrig, fileBad)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Feedback for the user running it that "it's working"
	fmt.Println("Successfully opened file")

	// Defer closing until the end of the function (almost... you'll see below)
	defer xmlFile.Close()

	// Read in the xml doc
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// Make a place for the data
	var records Records

	// Now, if it fails to unmarshal, we know the XML isn't going to load later
	xerr := xml.Unmarshal(byteValue, &records)
	if xerr != nil {
		// os will care that this program still has the file open, so close it
		xmlFile.Close()
		// And then make it BAD
		err := os.Rename(fileOrig, fileBad)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(xerr)
	}

	// Add/update WK_DAY to match Day
	for i := 0; i < len(records.Record); i++ {
		// But don't do it if it's been done, right?
		if records.Record[i].Day != "" {
			records.Record[i].WK_DAY = records.Record[i].Day
			records.Record[i].DAY = records.Record[i].Day
			records.Record[i].Day = ""
		}
	}

	// Marshal it to write it out
	modified, err := xml.MarshalIndent(&records, "", "	")
	if err != nil {
		log.Fatal(err)
	}
	// Don't forget the prolog!
	modified = []byte(xml.Header + string(modified))

	// Close it up (rmemeber up above, I said we may close before the end of the func)
	xmlFile.Close()

	// And replace the existing file with a new one
	errw := ioutil.WriteFile(fileOrig, modified, 0666)
	if errw != nil {
		log.Fatal(errw)
	}

}
