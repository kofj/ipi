package handlers

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"github.com/sirupsen/logrus"
	"kofj.net/toolz/ip/pkg/ipdb"
	"kofj.net/toolz/ip/pkg/models"
)

var indexTpl = "index.tmpl"
var ErrorNotFoundOrInvlid = errors.New("not found or invalid ip")

func IpiPage(ctx *gin.Context) {
	var format = "html"
	var ua = ctx.Request.UserAgent()
	uaInfo := useragent.Parse(ua)
	if uaInfo.IsUnknown() {
		format = "json"
		logrus.WithField("ua", ua).Warn("unknown user agent")
	}

	var ip = ctx.Param("ip")
	if ip == "" {
		ip = ctx.ClientIP()
	}
	city, err := ipdb.City(ip)
	if err != nil {
		errResp(ctx, ip, format, err)
		logrus.WithError(err).Error("ipdb query failed")
		return
	}

	if city == nil {
		err = ErrorNotFoundOrInvlid
		errResp(ctx, ip, format, err)
		logrus.WithError(err).Error("ipdb query failed")
		return
	}

	asn, err := ipdb.ASN(ip)
	if err != nil {
		errResp(ctx, ip, format, err)
		logrus.WithError(err).Error("ipdb query failed")
		return
	}
	if asn == nil {
		err = ErrorNotFoundOrInvlid
		errResp(ctx, ip, format, err)
		logrus.WithError(err).Error("ipdb query failed")
		return
	}
	var asnInfo = &models.ASN{
		Number:       asn.AutonomousSystemNumber,
		Organization: asn.AutonomousSystemOrganization,
	}
	logrus.WithField("asn", asn).Warn("asn")

	var locale = "en"
	var location *models.Location
	logrus.WithField("locale", city.Location).Warn("locale")
	if city.Location.Latitude != 0 || city.Location.Longitude != 0 {
		location = &models.Location{
			City:      city.City.Names[locale],
			Country:   city.Country.Names[locale],
			Timezone:  city.Location.TimeZone,
			Latitude:  city.Location.Latitude,
			Longitude: city.Location.Longitude,
		}
	}
	var info = models.Info{
		IP:       ip,
		Locale:   locale,
		ASN:      asnInfo,
		Location: location,
	}

	ipResp(ctx, ip, format, info)
}

func ipResp(ctx *gin.Context, ip, format string, data models.Info) {
	switch format {
	case "json":
		ctx.JSON(200, data)
	default:
		infoJson, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			errResp(ctx, ip, format, err)
			logrus.WithError(err).Error("json marshal failed")
		}

		ctx.HTML(200, indexTpl, gin.H{
			"ip":   ip,
			"info": string(infoJson),
		})
	}
}
func errResp(ctx *gin.Context, ip, format string, err error) {
	switch format {
	case "json":
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
	default:
		ctx.HTML(200, indexTpl, gin.H{
			"ip":    ip,
			"error": err.Error(),
		})
	}
}
