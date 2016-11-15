package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mhenderson-so/godnsmadeeasy/src/godnsmadeeasy"
)

func main() {

	DMEClient, err := GoDNSMadeEasy.NewGoDNSMadeEasy(&GoDNSMadeEasy.GoDNSMadeEasy{
		APIKey:               "",
		SecretKey:            "",
		APIUrl:               "https://api.sandbox.dnsmadeeasy.com/V2.0/",
		DisableSSLValidation: true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	allDomains, err := DMEClient.ExportAllDomains()
	if err != nil {
		fmt.Println(err)
		return
	}
	json, err := json.Marshal(allDomains)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(json))

	for domain, data := range *allDomains {
		fmt.Println(domain, data.Info.ID)
		timeStamp := time.Now().UnixNano()
		newRecord := &GoDNSMadeEasy.Record{
			Type:        "A",
			Name:        fmt.Sprintf("testrecord-%v", timeStamp),
			Value:       "127.0.0.1",
			GtdLocation: "DEFAULT",
			TTL:         300,
		}
		returnedRecord, err := DMEClient.AddRecord(data.Info.ID, newRecord)
		if err != nil {
			fmt.Println(err)
			return
		}
		if returnedRecord.Failed {
			fmt.Printf("Create record failed:%v+\n", returnedRecord)
		} else {
			fmt.Println("-- Created record OK:", returnedRecord.ID, "--")
		}

		returnedRecord.Name = fmt.Sprintf("postupdate-%s", returnedRecord.Name)
		DMEClient.UpdateRecord(data.Info.ID, returnedRecord)
		if err != nil {
			fmt.Println(err)
			return
		}
		if returnedRecord.Failed {
			fmt.Printf("Update record failed:%v+\n", returnedRecord)
		} else {
			fmt.Println("-- Update record OK:", returnedRecord.ID, "--")
		}

		err = DMEClient.DeleteRecord(data.Info.ID, returnedRecord.ID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("-- Deleted record OK:", returnedRecord.ID, "--")

	}

	/*
		fmt.Println("-- Getting domain list")
		domains, err := DMEClient.Domains()
		if err != nil {
			fmt.Println(err)
			return
		}

			domainID := domains[len(domains)-1].ID
			fmt.Println("-- Domain", domainID)

			domain, err := DMEClient.Domain(domainID)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%+v\n", domain)

			fmt.Println("-- Domain Records for", domainID)
			records, err := DMEClient.Records(domainID)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%+v\n", records)

				recordID := records[len(records)-1].ID
				fmt.Println("-- Domain Record for", domainID, ",", recordID)
				record, err := DMEClient.Record(domainID, recordID)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("%+v\n", record)


			fmt.Println("-- SOA")
			soa, err := DMEClient.SOA()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%+v\n", soa)

			fmt.Println("-- Vanity Nameservers")
			vanity, err := DMEClient.Vanity()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%+v\n", vanity)
	*/
}
