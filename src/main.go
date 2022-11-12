package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/dns/v2/recordsets"
	"github.com/rackspace/gophercloud/pagination"
)

type RecordSet struct {
	Name    string
	Type    string
	TTL     int
	Records []string
}

func main() {
	domainName := flag.String("domain", "", "Domain name to query")
	recordType := flag.String("type", "", "Record type to query")
	recordName := flag.String("name", "", "Record name to query")
	tenantName := flag.String("tenant", "", "Tenant name to query")
	authUrl := flag.String("url", "", "Auth URL to query")
	username := flag.String("username", "", "Username to query")
	password := flag.String("password", "", "Password to query")
	flag.Parse()

	if *domainName == "" || *recordType == "" || *recordName == "" || *tenantName == "" || *authUrl == "" || *username == "" || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	provider, err := openstack.AuthenticatedClient(gophercloud.AuthOptions{
		IdentityEndpoint: *authUrl,
		Username:         *username,
		Password:         *password,
		TenantName:       *tenantName,
	})
	if err != nil {
		panic(err)
	}

	dnsClient, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		panic(err)
	}

	var recordSet RecordSet
	var recordSets []RecordSet

	pager := recordsets.List(dnsClient, *domainName, recordsets.ListOpts{Name: *recordName, Type: *recordType})
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		recordSetList, err := recordsets.ExtractRecordSets(page)
		if err != nil {
			return false, err
		}

		for _, recordSet := range recordSetList {
			recordSet := RecordSet{
				Name:    recordSet.Name,
                Type:    recordSet.Type,
                TTL:     recordSet.TTL,
                Records: recordSet.Records,
            }
            recordSets = append(recordSets, recordSet)
        }
        return true, nil
    })
    if err != nil {
        panic(err)
    }

    for _, recordSet := range recordSets {
        fmt.Printf("%s %d %s %s", recordSet.Name, recordSet.TTL, recordSet.Type, strings.Join(recordSet.Records, " "))
    }
}
