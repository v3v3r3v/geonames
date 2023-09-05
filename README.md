# Geonames

#### Golang parsing library for the [geonames.org](http://www.geonames.org) database [dump](http://download.geonames.org/export/dump/).

![](static/example.gif)

## Features
- Parse data directly without downloading and unzipping
- Read line by line with low memory consumption

## Implemented data

|status|archive|comment|
|---|---|---|
|✅|xx.zip|GetGeonames; See [readme](#parsing-alphabetical-list-of-archives)|
|✅|admin1CodesASCII.txt|GetAdminDivisions|
|✅|admin2Codes.txt|GetAdminSubdivisions|
|✅|adminCode5.zip|GetAdminCodes5|
|✅|allCountries.zip|GetGeonames|
|⚠️|alternateNames.zip|depricated, use alternateNamesV2.zip instead|
|✅|alternateNamesDeletes-xxxx-xx-xx.txt|GetAlternateNameDeletes|
|✅|alternateNamesModifications-xxxx-xx-xx.txt|GetAlternateNameModifications|
|✅|alternateNamesV2.zip|GetAlternateNames|
|✅|alternatenames/xx.zip|GetAlternateNames; See [readme](#parsing-alphabetical-list-of-archives)|
|✅|cities1000.zip|GetGeonames|
|✅|cities15000.zip|GetGeonames|
|✅|cities500.zip|GetGeonames|
|✅|cities5000.zip|GetGeonames|
|✅|countryInfo.txt|GetCountries|
|✅|deletes-xxxx-xx-xx.txt|GetDeletes|
|✅|featureCodes_bg.txt|GetFeatureCodes|
|✅|featureCodes_en.txt|GetFeatureCodes|
|✅|featureCodes_nb.txt|GetFeatureCodes|
|✅|featureCodes_nn.txt|GetFeatureCodes|
|✅|featureCodes_no.txt|GetFeatureCodes|
|✅|featureCodes_ru.txt|GetFeatureCodes|
|✅|featureCodes_sv.txt|GetFeatureCodes|
|✅|hierarchy.zip|GetHierarchy|
|✅|iso-languagecodes.txt|GetLanguages|
|✅|modifications-xxxx-xx-xx.txt|GetModifications|
|✅|no-country.zip|GetGeonames|
|✅|shapes_all_low.zip|GetShapes|
|⚠️|shapes_simplified_low.json.zip|I don't see the point in geojson parsing|
|✅|timeZones.txt|GetTimeZones|
|✅|userTags.zip|GetUserTags|

## Installation

```bash 
$ go get github.com/v3v3r3v/geonames
```

## Quick start

#### Parsing cities
```go

package main

import (
    "fmt"
    "github.com/v3v3r3v/geonames"
    "github.com/v3v3r3v/geonames/models"
    "log"
)

func main() {
    p := geonames.NewParser()
    
    //print all cities with a population greater than 5000
    err := p.GetGeonames(geonames.Cities5000, func(geoname *models.Geoname) error {
    fmt.Println(geoname.Name)
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
}

```
#### Parsing alternames

```go
package main

import (
    "fmt"
    "github.com/v3v3r3v/geonames"
    "github.com/v3v3r3v/geonames/models"
    "log"
)
func main() {
    p := geonames.NewParser()
    
    err := p.GetAlternateNames(geonames.AlternateNames, func(geoname *models.AlternateName) error {
        fmt.Println(geoname.Name)
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

#### Parsing alphabetical list of archives

```go
package main

import (
    "fmt"
    "github.com/v3v3r3v/geonames"
    "github.com/v3v3r3v/geonames/models"
    "log"
)
func main() {
    p := geonames.NewParser()
    
    err := p.GetGeonames("AD.zip", func(geoname *models.Geoname) error {
        fmt.Println(geoname.Name)
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
    
    err = p.GetAlternateNames("alternames/AD.zip", func(geoname *models.AlternateName) error {
        fmt.Println(geoname.Name)
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
}
```
