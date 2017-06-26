package control

import (
	//"fmt"
  "log"
  "github.com/oschwald/geoip2-golang"
  "github.com/penolver/gowas/config"
  "net/http"
  "net"
)

var GeoIPDBdb *geoip2.Reader
var UseGeo bool
var Conf config.Configuration

func Validate(req *http.Request) bool {

  var returnstatus bool

  ipa, _, err := net.SplitHostPort(req.RemoteAddr)
  if err != nil { log.Println("userip: %q is not IP:port", req.RemoteAddr)}

  ip := net.ParseIP(ipa)
  // geo block
  if UseGeo == true {
    record, err := GeoIPDBdb.Country(ip)
    if err != nil {
      log.Println(err)
      //cityname = "unknown"
    } else {

      // if default allow, just block specific countries..
      if Conf.GeoIPDefaultAllow == true {
        returnstatus = true
        if _, ok := Conf.GeoIPDeniedCountries[record.Country.IsoCode]; ok { returnstatus = false }
        // else block all countries except certain allowed
      }else {
        returnstatus = false
        if _, ok := Conf.GeoIPAllowedCountries[record.Country.IsoCode]; ok { returnstatus = true }
      }
      //cityname = record.City.Names["en"]+`/`+record.Country.IsoCode
      log.Println("country is:",record.Country.IsoCode)

      //returnstatus = true
    }
  }
  // block localhost for testing specific blocks
  //if req.RemoteAddr == "127.0.0.1" {
    //returnstatus = false
  //}

  if returnstatus == true {
    return true
  }else { return false }
}
