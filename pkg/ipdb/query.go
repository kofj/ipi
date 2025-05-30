package ipdb

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/kofj/ipi/pkg/common"
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

func City(ctx context.Context, ip string) (city *geoip2.City, err error) {
	_, span := common.Tracer.Start(ctx, "ipi:ipdb:GeoCity")
	defer span.End()

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

func ASN(ctx context.Context, ip string) (asn *geoip2.ASN, err error) {
	_, span := common.Tracer.Start(ctx, "ipi:ipdb:GeoASN")
	defer span.End()

	if asnReader == nil {
		err = errors.New("IPDB not initialized")
		return
	}
	var ipNet = net.ParseIP(ip)
	if ipNet == nil {
		err = errors.New("invalid IP address")
		return
	}

	return asnReader.ASN(ipNet)
}
