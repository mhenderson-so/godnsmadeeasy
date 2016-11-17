# GoDNSMadeEasy
GoDNSMadeEasy is a GoLang library for accessing the DNS Made Easy API. It is not currently feature-complete. API coverage is as follows:

| Feature            | List | Create  | Update  | Delete 
|--------------------|---|---|---|---|
| Primary Domains    | ✓ | ✓ |  ✓ | ✓ 
| Secondary Domains  |   |   |   |   
| Records            | ✓ | ✓ | ✓ | ✓ 
| Vanity NS          | ✓ | ✓ | ✓ | ✓ 
| Custom SOA         | ✓ | ✓ | ✓ | ✓ 
| Templates          |   |   |   |   
| Transfer ACLs      |   |   |   |   
| Folders            |   |   |   |   
| Usage              |   |N/A| N/A  | N/A  
| Failover Monitor   |   |   |   |   N/A   
| IPset Fields       |   |   |   |    

# Usage

## GoDoc
Documentation for the API can be found at [https://godoc.org/github.com/mhenderson-so/godnsmadeeasy/src/GoDNSMadeEasy](https://godoc.org/github.com/mhenderson-so/godnsmadeeasy/src/GoDNSMadeEasy)

## Basic Usage
Create a new DNS Made Easy client with `NewGoDNSMadeEasy`:
```Go
import "github.com/mhenderson-so/godnsmadeeasy/src/GoDNSMadeEasy"

DMEClient, err := GoDNSMadeEasy.NewGoDNSMadeEasy(&GoDNSMadeEasy.GoDMEConfig{
    APIKey:               "d775b7a7-8192-46d2-80e8-53b95fda4931",
    SecretKey:            "c69f34e9-d8bc-4e0d-99b6-59476e73b61d",
    APIUrl:               GoDNSMadeEasy.LIVEAPI,
    DisableSSLValidation: false,
})
if err != nil {
    panic(err)
}

```
You can then retrieve data (domains, records, SOA, Nameservers) and update data (Create new domains, records, etc)

```Go
newDomain, err := DMEClient.AddDomain(&GoDNSMadeEasy.Domain{
    Name: "example.org",
})
if err != nil {
    panic(err)
}

newRecord, err := DMEClient.AddRecord(newDomain.ID, &GoDNSMadeEasy.Record{
    Type:        "A",
    Name:        "admin",
    Value:       "172.17.5.4",
    GtdLocation: "DEFAULT",
    TTL:         1800,
})
if err != nil {
    panic(err)
}
```

## Sample Application

There is a tiny sample application that is in the root folder of this project. This application just takes
your DNS Made Easy domain configuration and dumps it to `stdout` in JSON format. It's not intended to be
very useful - rather just as an example of how to use the library.

```
go run .\main.go -APIKey d775b7a7-8192-46d2-80e8-53b95fda4931 -SecretKey c69f34e9-d8bc-4e0d-99b6-59476e73b61d

{
	"example.org": {
		"SOA": null,
		"Info": {
			"name": "example.com",
			"id": 654321,
			"folderId": 1337,
			"nameServers": null,
			"updated": 1479328922220,
			"created": 1479254400000
		},
		"DefaultNS": null,
		"Records": []
	}
}
```

# Alternatives

This is far from the only DNS Made Easy integration built in Go. Here are some of the others,
and why I chose to write my own. I won't be offended if you choose one of the others over this library:

- [https://github.com/huguesalary/dnsmadeeasy](https://github.com/huguesalary/dnsmadeeasy): Not production ready
- [https://github.com/soniah/dnsmadeeas](https://github.com/soniah/dnsmadeeasy): Doesn't have full coverage of the API (neither does this API, yet) 
- [https://github.com/RealLeo/dnsmadeeasy](https://github.com/RealLeo/dnsmadeeasy): Doesn't have full coverage of the API (neither does this API, yet)
- [https://github.com/jswank/dnsme](https://github.com/jswank/dnsme): Complete coverage of the API, but is not a library, not easily convertable to a library, and uses the old v1.2 API