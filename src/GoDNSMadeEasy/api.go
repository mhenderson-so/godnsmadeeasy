package GoDNSMadeEasy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Domains returns the list of domains managed by DNS Made Easy
func (dme *GoDMEConfig) Domains() ([]Domain, error) {
	req, err := dme.newRequest("GET", "dns/managed/", nil)
	if err != nil {
		return nil, err
	}

	genericResponse := &GenericResponse{}
	err = dme.doDMERequest(req, &genericResponse)
	if err != nil {
		return nil, err
	}

	domainData := []Domain{}
	err = json.Unmarshal(genericResponse.Data, &domainData)
	return domainData, err
}

// Domain returns the summary data for a single domain. This is essentially the same as Domains(), but only returns one domain.
func (dme *GoDMEConfig) Domain(DomainID int) (*Domain, error) {
	reqStub := fmt.Sprintf("dns/managed/%v", DomainID)
	req, err := dme.newRequest("GET", reqStub, nil)
	if err != nil {
		return nil, err
	}

	domainResponse := &Domain{}
	err = dme.doDMERequest(req, domainResponse)
	if err != nil {
		return nil, err
	}

	return domainResponse, nil
}

// Records returns the records for a given domain. The domain is specified by its ID, which can be retrieved from Domains()
func (dme *GoDMEConfig) Records(DomainID int) ([]Record, error) {
	reqStub := fmt.Sprintf("dns/managed/%v/records", DomainID)
	req, err := dme.newRequest("GET", reqStub, nil)
	if err != nil {
		return nil, err
	}

	genericResponse := &GenericResponse{}
	err = dme.doDMERequest(req, &genericResponse)
	if err != nil {
		return nil, err
	}

	recordData := []Record{}
	err = json.Unmarshal(genericResponse.Data, &recordData)
	return recordData, err

}

// Record returns the record for a given record ID. This is essentially the same as Records(), but only returns one record
func (dme *GoDMEConfig) Record(DomainID, RecordID int) (*Record, error) {
	return nil, fmt.Errorf("Record() Not yet implemented")
}

// SOA returns custom Start of Authority records for an account.
func (dme *GoDMEConfig) SOA() ([]SOA, error) {
	req, err := dme.newRequest("GET", "dns/soa", nil)
	if err != nil {
		return nil, err
	}

	genericResponse := &GenericResponse{}
	err = dme.doDMERequest(req, &genericResponse)
	if err != nil {
		return nil, err
	}

	soaData := []SOA{}
	json.Unmarshal(genericResponse.Data, &soaData)

	return soaData, nil
}

// Vanity returns custom Vanity name servers for an account
func (dme *GoDMEConfig) Vanity() ([]Vanity, error) {
	req, err := dme.newRequest("GET", "dns/vanity", nil)
	if err != nil {
		return nil, err
	}

	vanityResponse := &GenericResponse{}
	err = dme.doDMERequest(req, &vanityResponse)
	if err != nil {
		return nil, err
	}

	vanityData := []Vanity{}
	json.Unmarshal(vanityResponse.Data, &vanityData)

	return vanityData, nil
}

// ExportAllDomains returns a map with every domain that DNS Made Easy manages, along with its properties
func (dme *GoDMEConfig) ExportAllDomains() (*AllDomainExport, error) {
	allDomains, err := dme.Domains()
	if err != nil {
		return nil, err
	}
	allSOA, err := dme.SOA()
	if err != nil {
		return nil, err
	}
	allVanity, err := dme.Vanity()
	if err != nil {
		return nil, err
	}

	thisExport := make(AllDomainExport)

	for _, domain := range allDomains {
		var thisSOA *SOA
		var thisVanity *Vanity

		//Find the correct SOA record
		for _, s := range allSOA {
			if s.ID == domain.SoaID {
				thisSOA = &s
			}
		}

		//Find the correct NS records
		for _, v := range allVanity {
			if v.ID == domain.VanityID {
				thisVanity = &v
			}
		}

		//Get DNS records
		thisRecords, err := dme.Records(domain.ID)
		if err != nil {
			return nil, err
		}

		thisExport[domain.Name] = DomainExport{
			Info:      &domain,
			SOA:       thisSOA,
			DefaultNS: thisVanity,
			Records:   &thisRecords,
		}
	}

	return &thisExport, nil
}

// AddRecord adds a DNS record to a given domain (identified by its ID)
func (dme *GoDMEConfig) AddRecord(DomainID int, RecordRecord *Record) (*Record, error) {
	reqStub := fmt.Sprintf("dns/managed/%v/records", DomainID)
	bodyData, err := json.Marshal(RecordRecord)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewReader(bodyData)
	req, err := dme.newRequest("POST", reqStub, bodyBuffer)
	if err != nil {
		return nil, err
	}

	returnedRecord := &Record{}
	err = dme.doDMERequest(req, returnedRecord)
	if err != nil {
		return nil, err
	}

	return returnedRecord, err
}

// UpdateRecord updates an existing DNS record (identified by its ID) in a given domain. DNS Made Easy only returns success/fail for this method.
func (dme *GoDMEConfig) UpdateRecord(DomainID int, Record *Record) error {
	reqStub := fmt.Sprintf("dns/managed/%v/records/%v", DomainID, Record.ID)
	bodyData, err := json.Marshal(Record)
	if err != nil {
		return err
	}
	bodyBuffer := bytes.NewReader(bodyData)
	req, err := dme.newRequest("PUT", reqStub, bodyBuffer)
	if err != nil {
		return err
	}

	//A PUT request does not return a body, so we'll do another Get
	err = dme.doDMERequest(req, nil)
	if err != nil {
		return err
	}

	return err
}

// DeleteRecord deletes an existing DNS record (identified by its ID) in a given domain
func (dme *GoDMEConfig) DeleteRecord(DomainID, RecordID int) error {
	reqStub := fmt.Sprintf("dns/managed/%v/records/%v", DomainID, RecordID)
	req, err := dme.newRequest("DELETE", reqStub, nil)
	if err != nil {
		return err
	}
	return dme.doDMERequest(req, nil)
}

// AddDomain adds a domain to your DNS Made Easy account
func (dme *GoDMEConfig) AddDomain(DomainRecord *Domain) (*Domain, error) {
	reqStub := "dns/managed/"
	bodyData, err := json.Marshal(DomainRecord)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewReader(bodyData)
	req, err := dme.newRequest("POST", reqStub, bodyBuffer)
	if err != nil {

		return nil, err
	}

	returnedDomain := &Domain{}
	err = dme.doDMERequest(req, returnedDomain)
	if err != nil {
		return nil, err
	}

	return returnedDomain, err
}

// DeleteDomain deletes a domain to your DNS Made Easy account. The DeleteTimeout argument indicates how long we should keep trying to
// delete the domain if DNS Made Easy says it can't delete the domain due to a pending operation. In these cases, usually deleting a domain
// name will succeed after a certain period of time. You may not want to wait for this time though, so specify 0 here to never retry.
func (dme *GoDMEConfig) DeleteDomain(DomainID int, DeleteTimeout time.Duration) error {
	timeOutAt := time.Now().Add(DeleteTimeout)

	reqStub := fmt.Sprintf("dns/managed/%v", DomainID)
	req, err := dme.newRequest("DELETE", reqStub, nil)
	if err != nil {
		return err
	}
	//Try to delete once
	deleteError := dme.doDMERequest(req, nil)
	if deleteError == nil || DeleteTimeout == 0 {
		return deleteError
	}

	//If we were unsuccessful in deleting the first time, try try again until the timeout
	for time.Now().Before(timeOutAt) {
		req, _ := dme.newRequest("DELETE", reqStub, nil)
		deleteError := dme.doDMERequest(req, nil)
		//No error? Then we're all done.
		if deleteError == nil {
			return deleteError
		}
		//We got a different error this time that is not a pending delete error
		if deleteError.Error() != pendingDeleteError {
			return deleteError
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("Could not delete domain after %s (%s)", DeleteTimeout.String(), err)

}
