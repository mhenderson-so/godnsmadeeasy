package GoDNSMadeEasy

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	apiKey         = flag.String("APIKey", "", "Your DNS Made Easy Sandbox API Key")
	secretKey      = flag.String("SecretKey", "", "Your DNS Made Easy Sandbox Secret Key")
	timeAdjust     = flag.Int("TimeOffset", 0, "Timestamp adjustment in seconds. DNS Made Easy has a very strict time synchronisation requirement. If your local clock runs slightly fast or slow (even by 30 seconds), requests will fail. You can adjust the timestamp sent by DNS Made Easy here to account for this offset")
	DomainsCreated = make(map[string]*Domain)
)

func TestMain(m *testing.M) {
	flag.Parse()
	m.Run()
	cleanUpDomains()
}

// TestCreateDomain tests the creation of a domain. This is kind of a redundant test, because every other test is going to fail
// if we can't do this.
func TestCreateDomain(t *testing.T) {
	DMEClient, err := newClient()
	if err != nil {
		t.Fatal(err)
	}
	newDomain, err := generateTestDomain(DMEClient)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Using test domain name", newDomain.Name)

	if newDomain.ID == 0 {
		t.Fatal("domain ID is 0")
	}
}

// TestCreateDomain tests the creation of a domain. This is kind of a redundant test, because every other test is going to fail
// if we can't do this.
func TestRecords(t *testing.T) {
	var TestRecords = getTestRecords(false)
	var UpdateRecords = getTestRecords(true)
	var CreatedRecords []*Record

	DMEClient, err := newClient()
	if err != nil {
		t.Fatal(err)
	}
	newDomain, err := generateTestDomain(DMEClient)
	if err != nil {
		t.Fatal(err)
	}
	DomainID := newDomain.ID
	t.Log("Using test domain name", newDomain.Name)

	//Create a record of each type
	for _, thisRecord := range TestRecords {
		newRecord, err := DMEClient.AddRecord(DomainID, &thisRecord)
		if err != nil {
			t.Error(fmt.Sprintf("%s: %s", thisRecord.Name, err))
		}
		mismatches := compareRecords(&thisRecord, newRecord)
		if len(mismatches) > 0 {
			t.Error(fmt.Sprintf("(create) %s %s: records do not match: %s", thisRecord.Type, thisRecord.Name, strings.Join(mismatches, ",")))
		}
		CreatedRecords = append(CreatedRecords, newRecord)
	}

	//Update previously created records
	for _, thisRecord := range UpdateRecords {
		for _, existingRecord := range CreatedRecords {
			if thisRecord.Type == existingRecord.Type && thisRecord.Name == existingRecord.Name {
				thisRecord.ID = existingRecord.ID
				err := DMEClient.UpdateRecord(DomainID, &thisRecord)
				if err != nil {
					t.Error(fmt.Sprintf("(update) %s %s: %s", thisRecord.Type, thisRecord.Name, err))
				}
				//Because DNS Made Easy does not return the new record, and doesn't give a method for retrieving a single record, this is
				//all we can do here
			}
		}
	}

	//And delete ther records we just Updated
	for _, existingRecord := range CreatedRecords {
		err := DMEClient.DeleteRecord(DomainID, existingRecord.ID)
		if err != nil {
			t.Error(fmt.Sprintf("(delete) %s %s: %s", existingRecord.Type, existingRecord.Name, err))
		}
	}
}

//Create a test domain, return the domain entry for this domain, and add it to our list of domains that needs to be cleaned up at the end
//Names are generated using a timestamp.
func generateTestDomain(DMEClient *GoDMEConfig) (*Domain, error) {
	thisDomainName := fmt.Sprintf("gotest-%v.org", time.Now().UnixNano())
	newDomain, err := DMEClient.AddDomain(&Domain{
		Name: thisDomainName,
	})
	if err != nil {
		return nil, err
	}

	DomainsCreated[thisDomainName] = newDomain
	return newDomain, nil
}

//We need to clean up after our tests are run, so we don't leave old domains lying around in the sandbox
func cleanUpDomains() {
	fmt.Println("Cleaning up domains...")
	//Create a client for talking to DME
	DMEClient, err := newClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	//Create a WaitGroup, so we can delete the domains in parallel, but wait for all to complete
	var wg sync.WaitGroup
	for name, domain := range DomainsCreated { //Loop through the domains we created during this testing
		wg.Add(1)                              //Add one to the wait group
		go func(name string, domain *Domain) { //Delete the domains asynchronously
			defer wg.Done()                                         //When this is finished, indicate to the Wait Group that we're done
			fmt.Println("Deleting", name)                           //Send something to console so we know what's going on
			err := DMEClient.DeleteDomain(domain.ID, 2*time.Minute) //Delete the domain, with a 2 minute timeout. Sandbox takes around 50 seconds on average
			if err != nil {
				fmt.Println("Could not delete", name, "error:", err)
			}
		}(name, domain)
	}
	wg.Wait() //Wait for all the Done()'s to come through
}

//Create a DNS Made Easy client for each test to run from, as they are run in parallel
func newClient() (*GoDMEConfig, error) {
	return NewGoDNSMadeEasy(&GoDMEConfig{
		APIKey:               *apiKey,
		SecretKey:            *secretKey,
		APIUrl:               SANDBOXAPI,
		DisableSSLValidation: true,
		TimeAdjust:           (time.Duration(*timeAdjust) * time.Second),
	})

}

func getTestRecords(Updated bool) []Record {
	recIPVal, recTTL, recIPv6Val, recDomain, recData := "127.8.4.3", 300, "::1", "example.org.", "\"originalvalue\""

	if Updated {
		recIPVal, recTTL, recIPv6Val, recDomain, recData = "10.85.67.244", 1800, "::BEEF", "example.com.", "\"newvalue\""
	}
	var TestRecords []Record

	//Gimmie an A
	TestRecords = append(TestRecords, Record{
		Name:        "testa",
		Type:        "A",
		Value:       recIPVal,
		TTL:         recTTL,
		GtdLocation: "DEFAULT",
	})

	//Gimmie an AAAA
	TestRecords = append(TestRecords, Record{
		Name:        "testaaaa",
		Type:        "AAAA",
		Value:       recIPv6Val,
		TTL:         recTTL,
		GtdLocation: "DEFAULT",
	})

	//Gimmie a CNAME
	TestRecords = append(TestRecords, Record{
		Name:        "testcname",
		Type:        "CNAME",
		Value:       recDomain,
		TTL:         recTTL,
		GtdLocation: "DEFAULT",
	})

	//Gimmie a ANAME
	TestRecords = append(TestRecords, Record{
		Name:        "",
		Type:        "ANAME",
		Value:       recDomain,
		TTL:         recTTL,
		GtdLocation: "DEFAULT",
	})

	//Gimmie a MX
	TestRecords = append(TestRecords, Record{
		Name:        "testmx",
		Type:        "MX",
		Value:       recDomain,
		TTL:         recTTL,
		MxLevel:     10,
		GtdLocation: "DEFAULT",
	})

	//Gimmie a HTTP
	TestRecords = append(TestRecords, Record{
		Name:         "testred",
		Type:         "HTTPRED",
		Value:        strings.TrimSuffix(fmt.Sprintf("http://%s", recDomain), "."),
		TTL:          recTTL,
		HardLink:     false,
		RedirectType: "STANDARD - 301",
		Title:        "test redirect title",
		Keywords:     "just,stuff",
		Description:  "just doin some stuff",
		GtdLocation:  "DEFAULT",
	})

	//Gimmie a TXT
	TestRecords = append(TestRecords, Record{
		Name:        "testtxt",
		Type:        "TXT",
		Value:       recData,
		TTL:         recTTL,
		GtdLocation: "DEFAULT",
	})

	//Gimmie a SPF
	TestRecords = append(TestRecords, Record{
		Name:        "testtxt",
		Type:        "SPF",
		Value:       recData,
		TTL:         recTTL,
		GtdLocation: "DEFAULT",
	})

	//Gimmie a PTR. Yeah I know this isn't a useful PTR record, but we can still test with it
	TestRecords = append(TestRecords, Record{
		Name:        "testptr",
		Type:        "PTR",
		Value:       recDomain,
		TTL:         recTTL,
		GtdLocation: "DEFAULT",
	})

	//Gimmie a NS
	TestRecords = append(TestRecords, Record{
		Name:        "testns",
		Type:        "NS",
		Value:       recDomain,
		TTL:         recTTL,
		GtdLocation: "DEFAULT",
	})

	//Gimmie a SRV
	TestRecords = append(TestRecords, Record{
		Name:        "_testsrv",
		Type:        "SRV",
		Priority:    10,
		Weight:      10,
		Port:        80,
		Value:       recDomain,
		TTL:         recTTL,
		GtdLocation: "DEFAULT",
	})

	return TestRecords
}

func compareRecords(a, b *Record) []string {

	var mismatches []string
	if a == nil || b == nil {
		if a == nil {
			mismatches = append(mismatches, "A is nil")
		}

		if a == nil {
			mismatches = append(mismatches, "B is nil")
		}
		return mismatches
	}

	//All records have a name, a type, a value,  a TTL and a GtdLocation
	if a.Type != b.Type {
		mismatches = append(mismatches, "Type")
	}
	if a.Name != b.Name {
		mismatches = append(mismatches, "Name")
	}
	if a.Value != b.Value {
		mismatches = append(mismatches, "Value")
	}
	if a.TTL != b.TTL {
		mismatches = append(mismatches, "TTL")
	}
	if a.GtdLocation != b.GtdLocation {
		mismatches = append(mismatches, "GtdLocation")
	}

	//But some have more
	switch a.Type {
	case "MX":
		if a.MxLevel != b.MxLevel {
			mismatches = append(mismatches, "MxLevel")
		}

	case "HTTP":
		if a.HardLink != b.HardLink {
			mismatches = append(mismatches, "HardLink")
		}

		if a.Title != b.Title {
			mismatches = append(mismatches, "Title")
		}

		if a.Keywords != b.Keywords {
			mismatches = append(mismatches, "Keywords")
		}

		if a.Description != b.Description {
			mismatches = append(mismatches, "Description")
		}
	case "SRV":
		if a.Weight != b.Weight {
			mismatches = append(mismatches, "Weight")
		}
		if a.Port != b.Port {
			mismatches = append(mismatches, "Port")
		}
		if a.Priority != b.Priority {
			mismatches = append(mismatches, "Priority")
		}
	}
	return mismatches
}
