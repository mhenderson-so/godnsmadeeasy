package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/mhenderson-so/godnsmadeeasy/src/GoDNSMadeEasy"
)

var (
	apiKey     = flag.String("APIKey", "", "Your DNS Made Easy API Key")
	secretKey  = flag.String("SecretKey", "", "Your DNS Made Easy Secret Key")
	sandbox    = flag.Bool("Sandbox", false, "Use the DNS Made Easy Sandbox API")
	timeAdjust = flag.Int("TimeOffset", 0, "Timestamp adjustment in seconds. DNS Made Easy has a very strict time synchronisation requirement. If your local clock runs slightly fast or slow (even by 30 seconds), requests will fail. You can adjust the timestamp sent by DNS Made Easy here to account for this offset")
)

func main() {
	//Parse command-line flags
	flag.Parse()

	//Validate that the appropriate flags have been provided
	if *apiKey == "" {
		fmt.Println("You must provide your DNS Made Easy API Key")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		return
	}

	if *secretKey == "" {
		fmt.Println("You must provide your DNS Made Easy Secret Key")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		return
	}

	//Use the normal API unless we want to talk to the sandbox API
	apiURL := "https://api.dnsmadeeasy.com/V2.0/"
	if *sandbox {
		apiURL = "https://api.sandbox.dnsmadeeasy.com/V2.0/"
	}

	//Create our client for talking to DNS Made Easy, using the sandbox API
	DMEClient, err := GoDNSMadeEasy.NewGoDNSMadeEasy(&GoDNSMadeEasy.GoDNSMadeEasy{
		APIKey:               *apiKey,
		SecretKey:            *secretKey,
		APIUrl:               apiURL,
		DisableSSLValidation: *sandbox, //Only disable SSL validation for the sandbox
		TimeAdjust:           (time.Duration(*timeAdjust) * time.Second),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//For this demo, all we're going to do is invoke the ExportAllDomains() function, which gets all of our
	//domains from DNS Made Easy and dumps their data into a single object. This makes several read calls to the
	//API, and is good enough for testing connectivity and credentials to the API
	allDomains, err := DMEClient.ExportAllDomains()
	if err != nil {
		fmt.Println(err)
		return
	}

	//Pretty-print our domain data
	json, err := json.MarshalIndent(allDomains, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	//Dump the exported data as JSON, pretty-printed to the console
	fmt.Println(string(json))
}
