package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

var (
	appVersion = "1.0.1"
)

var (
	regionID        = ""
	accessKeyID     = ""
	accessKeySecret = ""
	domainName      = ""
	domainType      = ""
	rrName          = ""
)

// Env constants
const (
	EnvPrefix       = "ALIDNS_"
	EnvRegion       = EnvPrefix + "REGION"
	EnvAccessKey    = EnvPrefix + "ACCESS_KEY"
	EnvAccessSecret = EnvPrefix + "ACCESS_SECRET"
	EnvDomain       = EnvPrefix + "DOMAIN"
	EnvDomainType   = EnvPrefix + "DDMAIN_TYPE"
	EnvRR           = EnvPrefix + "RR"
)

// Default values
const (
	DefRegion = "cn-hangzhou"
	DefDomainType = "A"
)

func main() {
	parseArgs()
	publicIP, err := lookupPublicIP()
	if err != nil {
		fmt.Println("Loopup self IP failed:", err.Error())
		os.Exit(1)
	}
	dns, err := newAliDNSClient()
	if err != nil {
		fmt.Println("Create request client error", err.Error())
		os.Exit(1)
	}
	ip, recordID, err := fetchResolvingStatus(dns)
	if err != nil {
		fmt.Println("Fetch resolved ip error", err.Error())
		os.Exit(1)
	}
	if publicIP == ip {
		fmt.Println("IP no change")
	} else {
		fmt.Printf("My public IP now %s, resolving IP %s, updating...\n", publicIP, ip)
		err = updateRecord(dns, publicIP, recordID)
		if err != nil {
			fmt.Println("Update error", err.Error())
			os.Exit(1)
		}
		fmt.Println("Update successfully")
	}
}

func parseArgs() {
	flag.StringVar(&regionID, "region", DefRegion, "Region ID")
	flag.StringVar(&accessKeyID, "key", "", "Access Key ID")
	flag.StringVar(&accessKeySecret, "secret", "", "Access Key Secret")
	flag.StringVar(&domainName, "domain", "", "Domain name (like google.com)")
	flag.StringVar(&rrName, "rr", "", "Resource record (RR)")
	flag.StringVar(&domainType, "type", DefDomainType, "Domain type (A,CNAME,MX,etc...)")
	ver := flag.Bool("v", false, "Show version")
	flag.Parse()
	if *ver {
		fmt.Println("version", appVersion)
		os.Exit(0)
	}
	checkArg(&regionID, "Region ID", "region", EnvRegion, DefRegion)
	checkArg(&accessKeyID, "Access Key ID", "key", EnvAccessKey, "")
	checkArg(&accessKeySecret, "Access Key Secret", "secret", EnvAccessSecret, "")
	checkArg(&domainName, "Domain", "domain", EnvDomain, "")
	checkArg(&rrName, "Resource record (RR)", "rr", EnvRR, "")
	checkArg(&domainType, "Domain type", "type", EnvDomainType, DefDomainType)
}

func checkArg(v *string, name, para, env, def string) {
	if *v == "" {
		vv, exist := os.LookupEnv(env)
		if !exist {
			vv = def
			if vv == "" {
				fmt.Printf("%s is required, specify by -%s parameter or %s env\n", name, para, env)
				os.Exit(1)
			}
		}
		*v = vv
	}
}

func newAliDNSClient() (*alidns.Client, error) {
	c, err := alidns.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func fetchResolvingStatus(dns *alidns.Client) (ip, recordID string, err error) {
	req := alidns.CreateDescribeDomainRecordsRequest()
	req.DomainName = domainName
	req.RRKeyWord = rrName
	res, err := dns.DescribeDomainRecords(req)
	if err != nil {
		return "", "", err
	}
	if len(res.DomainRecords.Record) == 0 {
		return "", "", fmt.Errorf("Could not find record of %s.%s", rrName, domainName)
	}
	record := res.DomainRecords.Record[0]
	return record.Value, record.RecordId, nil
}

func updateRecord(dns *alidns.Client, ip, recordID string) error {
	req := alidns.CreateUpdateDomainRecordRequest()
	req.RecordId = recordID
	req.Type = "A"
	req.Value = ip
	req.RR = rrName
	_, err := dns.UpdateDomainRecord(req)
	if err != nil {
		return err
	}
	return nil
}

func lookupPublicIP() (string, error) {
	res, err := http.Get("http://whatismyip.akamai.com")
	if err != nil {
		return "", errors.New("Query my IP failed")
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("Query my IP failed")
	}
	myIP := string(b)
	resolvedIP, err := resolveClosestRouteIP(myIP)
	if err != nil {
		return "", err
	}
	if resolvedIP != myIP {
		return "", errors.New("Non public ip")
	}
	return myIP, nil
}
