package ipdb

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

var asnOnce sync.Once
var cityOnce sync.Once
var asnReader *geoip2.Reader
var cityReader *geoip2.Reader

// const databaseCN = "data/qqwry.mmdb"

const asnDb = "data/GeoLite2-ASN.mmdb"
const cityDb = "data/GeoLite2-City.mmdb"

// GetIPDBInstance returns a singleton instance of IPDB.
func Singleton() (err error) {
	cityOnce.Do(func() {
		cityReader, err = geoip2.Open(cityDb)
	})
	if err != nil {
		return
	}
	asnOnce.Do(func() {
		asnReader, err = geoip2.Open(asnDb)
	})
	return
}

func City(ip string) (city *geoip2.City, err error) {
	if cityReader == nil {
		err = errors.New("IPDB not initialized")
		return
	}
	var ipNet = net.ParseIP(ip)
	if ipNet == nil {
		err = errors.New("invalid IP address")
		return
	}
	return cityReader.City(ipNet)
}

func ASN(ip string) (asn *geoip2.ASN, err error) {
	if asnReader == nil {
		err = errors.New("IPDB not initialized")
		return
	}
	var ipNet = net.ParseIP(ip)
	if ipNet == nil {
		err = errors.New("invalid IP address")
		return
	}

	isp, err := asnReader.ISP(ipNet)
	fmt.Println(isp)
	fmt.Println(err)

	return asnReader.ASN(ipNet)
}
