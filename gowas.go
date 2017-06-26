package main

import (
  "net/http"
  "crypto/tls"
  "github.com/vulcand/oxy/forward"
  "github.com/vulcand/oxy/testutils"
  "github.com/oschwald/geoip2-golang"
  "github.com/penolver/gowas/config"
  "github.com/penolver/gowas/control"
  "log"
  "fmt"
  )

const version = "0.0.2"

func main() {

fmt.Println(`
              _       _____   _____
   ____ _____| |     / /   | / ___/
  / __ \/ __ \ | /| / / /| | \__ \
 / /_/ / /_/ / |/ |/ / ___ |___/ /
 \__, /\____/|__/|__/_/  |_/____/
/____/  goWAS version `+version+`
`)

  // load config..
  log.Println("Loading Config..")
  config := config.Load("./config.json")
  control.Conf = config

  var err error

  // if GeoIP DB provided in config
  log.Println("Loading GeoIP DB..")
  if config.GeoIPDB != "" {
    control.GeoIPDBdb, err = geoip2.Open(config.GeoIPDB)
  }else{
    control.GeoIPDBdb, err = geoip2.Open("GeoLite2-Country.mmdb")
  }
  if err != nil {
    log.Println("Issue with geo database..",err," continuing without..")
    control.UseGeo = false
  }else {
    control.UseGeo = true
    defer control.GeoIPDBdb.Close()
  }


  // Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
  fwd, _ := forward.New()

  redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
      // validate request
      check := control.Validate(req)
      if check == false {
        log.Println("Request denied from: ",req.RemoteAddr)
        //log.Println("full req details: ",req)
        w.Write([]byte("Request Denied\n"))
      } else {
        // let us forward this request to another server
        log.Println("Request reverse proxied")
    		req.URL = testutils.ParseURI(config.ForwardTo)
    		fwd.ServeHTTP(w, req)
      }

  })

  // if HTTPS Listener
  if config.ListenProto == "TLS" {
    cfg := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        },
    }
    srv := &http.Server{
        Addr:         ":"+config.ListenPort,
        Handler:      redirect,
        TLSConfig:    cfg,
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
    }
    log.Println("Starting HTTPS (TLS 1.2) server on :",config.ListenPort)
    log.Fatal(srv.ListenAndServeTLS(config.TLS_Cert, config.TLS_Key))

  // else assume HTTP listener
  } else {
    s := &http.Server{
    	Addr:           ":"+config.ListenPort,
    	Handler:        redirect,
    }
    log.Println("Starting HTTP server on: ",config.ListenPort)
    s.ListenAndServe()
  }

}
