package geonames

import (
	"geonames/models"
	"geonames/stream"
	"io"
	"net/http"
	"strings"
)

const Url = "https://download.geonames.org/export/dump/"

type GeoNameFile string
type AltNameFile string
type LangCodeFile string

//List of dump archives
const (
	Cities500      GeoNameFile  = "cities500.zip"
	Cities1000     GeoNameFile  = "cities1000.zip"
	Cities5000     GeoNameFile  = "cities5000.zip"
	Cities15000    GeoNameFile  = "cities15000.zip"
	AllCountries   GeoNameFile  = "allCountries.zip"
	NoCountry      GeoNameFile  = "no-country.zip"
	AlternateNames AltNameFile  = "alternateNamesV2.zip"
	LangCodes      LangCodeFile = "iso-languagecodes.txt"
)

type Parser func(file string) (io.ReadCloser, error)

func NewParser() Parser {
	return Parser(func(file string) (io.ReadCloser, error) {
		resp, err := http.Get(Url + file)
		if err != nil {
			return nil, err
		}

		return resp.Body, nil
	})
}

func (p Parser) GetGeonames(archive GeoNameFile, handler func(*models.Geoname) error) error {
	return p.get(string(archive), func(columns []string) error {
		model, err := models.ParseGeoname(columns)
		if err != nil {
			return err
		}

		return handler(model)
	})
}

func (p Parser) GetAlternames(handler func(*models.Altername) error) error {
	return p.get(string(AlternateNames), func(columns []string) error {
		model, err := models.ParseAltername(columns)
		if err != nil {
			return err
		}

		return handler(model)
	})
}

func (p Parser) get(archive string, handler func(columns []string) error) error {
	r, err := p(archive)
	if err != nil {
		return err
	}

	return stream.Stream(r, defaultFilename(archive), func(columns []string) error {
		return handler(columns)
	})
}

func defaultFilename(archive string) string {
	return strings.Replace(string(archive), ".zip", ".txt", 1)
}
