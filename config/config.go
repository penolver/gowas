// load JSON config

package config

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "os"
)

type Configuration struct {
  ListenPort  string `json:"ListenPort"`
  ListenProto string `json:"ListenProto"`
  TLS_Cert   string `json:"TLS_Cert"`
  TLS_Key string `json:"TLS_Key"`
  ForwardTo string `json:"ForwardTo"`
  GeoIPDB  string `json:"GeoIPDB"`
  GeoIPDefaultAllow bool `json:"GeoIPDefaultAllow"`
  // using a hash map for lookup speed, looks a bit odd in config admittedly, e.g. "GB": "GB"
  GeoIPAllowedCountries map[string]string `json:"GeoIPAllowedCountries"`
  GeoIPDeniedCountries map[string]string `json:"GeoIPDeniedCountries"`
  Verbose_Logging bool `json:"Verbose_Logging"`

}

// Load JSON Config file
func Load(path string) Configuration {
	file, err := ioutil.ReadFile(path)
	if err != nil {
    configDefault := `
    {
      "ListenPort":"8890",
      "ListenProto":"TLS",
      "TLS_Cert":"cert.pem",
      "TLS_Key":"key.pem",
      "ForwardTo":"https://atom.io:443",
      "GeoIPDB":"GeoLite2-Country.mmdb",
      "GeoIPDefaultAllow": true,
      "GeoIPAllowedCountries":[
        {"Country": "GB"},
        {"Country": "US"}
      ],
      "GeoIPDeniedCountries":[
        {"Country": "CN"},
        {"Country": "ALL"}
      ],
      "Verbose_Logging":true
    }`
    log.Println("Default config style (config.json):",configDefault)
		log.Fatal("Config File Missing. ", err)
	}

	var config Configuration
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}

	return config
} // LoadConfig

// Save JSON config file
func Save(v interface{}, path string) {
    fo, err := os.Create(path)
    if err != nil {
      log.Fatal("Config save file write error: ",err)
    }
    defer fo.Close()
    e := json.NewEncoder(fo)
    if err := e.Encode(v); err != nil {
      log.Fatal("Config save file encode error: ",err)
    }
} // SaveConfig
