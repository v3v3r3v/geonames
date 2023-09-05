package geonames

import (
	"errors"
	"fmt"
	"github.com/jszwec/csvutil"
	"github.com/v3v3r3v/geonames/models"
	"github.com/v3v3r3v/geonames/stream"
	"io"
	"net/http"
	"os"
)

const DownloadGeonamesOrgUrl = "https://download.geonames.org/export/dump/"

// List of dump archives
const (
	Cities500                   models.GeoNameFile     = "cities500.zip"
	Cities1000                  models.GeoNameFile     = "cities1000.zip"
	Cities5000                  models.GeoNameFile     = "cities5000.zip"
	Cities15000                 models.GeoNameFile     = "cities15000.zip"
	AllCountries                models.GeoNameFile     = "allCountries.zip"
	NoCountry                   models.GeoNameFile     = "no-country.zip"
	AlternateNames              models.AltNameFile     = "alternateNamesV2.zip"
	LangCodes                   models.DumpFile        = "iso-languagecodes.txt"
	TimeZones                   models.DumpFile        = "timeZones.txt"
	Countries                   models.DumpFile        = "countryInfo.txt"
	FeatureCodeBg               models.FeatureCodeFile = "featureCodes_bg.txt"
	FeatureCodeEn               models.FeatureCodeFile = "featureCodes_en.txt"
	FeatureCodeNb               models.FeatureCodeFile = "featureCodes_nb.txt"
	FeatureCodeNn               models.FeatureCodeFile = "featureCodes_nn.txt"
	FeatureCodeNo               models.FeatureCodeFile = "featureCodes_no.txt"
	FeatureCodeRu               models.FeatureCodeFile = "featureCodes_ru.txt"
	FeatureCodeSv               models.FeatureCodeFile = "featureCodes_sv.txt"
	Hierarchy                   models.DumpFile        = "hierarchy.zip"
	Shapes                      models.DumpFile        = "shapes_all_low.zip"
	UserTags                    models.DumpFile        = "userTags.zip"
	AdminDivisions              models.DumpFile        = "admin1CodesASCII.txt"
	AdminSubDivisions           models.DumpFile        = "admin2Codes.txt"
	AdminCode5                  models.DumpFile        = "adminCode5.zip"
	AlternateNamesDeletes       models.DumpFile        = "alternateNamesDeletes-%s.txt"
	AlternateNamesModifications models.DumpFile        = "alternateNamesModifications-%s.txt"
	Deletes                     models.DumpFile        = "deletes-%s.txt"
	Modifications               models.DumpFile        = "modifications-%s.txt"
)

type FetcherConfig struct {
	RemoteUrl string
	LocalPath string
}

type Fetcher struct {
	cfg FetcherConfig
}

func NewFetcher(cfg FetcherConfig) Fetcher {
	return Fetcher{
		cfg: cfg,
	}
}

type FetchSource string

const (
	SourceFs   = "fs"
	SourceHttp = "http"
)

func (p Fetcher) FetchFile(source FetchSource, file models.DumpFile) (io.ReadCloser, error) {
	switch source {
	case SourceFs:
		return p.fetchFs(file)
	case SourceHttp:
		return p.fetchHttp(file)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown fetch source: %s", source))
	}
}

func (p Fetcher) DumpToFile(file models.DumpFile) error {
	reader, err := p.fetchHttp(file)

	if err != nil {
		return err
	}

	defer reader.Close()

	// Create or open a file for writing
	fsFile, err := os.Create(p.cfg.LocalPath + file.String())
	if err != nil {
		return err
	}
	defer fsFile.Close()

	// Copy the response body to the file
	_, err = io.Copy(fsFile, reader)
	if err != nil {
		return err
	}

	return nil
}

func (p Fetcher) fetchFs(file models.DumpFile) (io.ReadCloser, error) {
	fsFile, err := os.Open(p.cfg.LocalPath + file.String())
	if err != nil {
		return nil, err
	}
	return io.NopCloser(fsFile), nil
}

func (p Fetcher) fetchHttp(file models.DumpFile) (io.ReadCloser, error) {
	url := p.cfg.RemoteUrl + file.String()
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Http Get %s error: %s", url, err))
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp.Body, nil
	case http.StatusNotFound:
		return nil, errors.New(fmt.Sprintf("Page %s does not exist", url))
	default:
		return nil, errors.New(fmt.Sprintf("Page %s returned unexpected code %d", url, resp.StatusCode))
	}
}

type Parser struct {
	FetchSource FetchSource
	Fetcher     Fetcher
}

func handle[T any](p Parser, dump models.DumpFile, isHeadersEmpty bool, handler func(*T) error) error {
	var err error
	var headers = []string{}

	model := new(T)

	if isHeadersEmpty {
		headers, err = csvutil.Header(model, "csv")
		if err != nil {
			return err
		}
	}

	r, err := p.Fetcher.FetchFile(p.FetchSource, dump)
	if err != nil {
		return err
	}

	f := func(parse func(v interface{}) error) error {
		v := new(T)
		err = parse(v)
		if err != nil {
			return err
		}
		return handler(v)
	}

	if dump.IsArchive() {
		err = stream.StreamArchive(r, dump.TextFilename(), f, headers)
	} else {
		err = stream.StreamFile(r, f, headers)
	}

	return err
}

func (p Parser) GetGeonames(archive models.GeoNameFile, handler func(*models.Geoname) error) error {
	return handle(p, archive.DumpFile(), true, handler)
}

func (p Parser) GetAlternateNames(archive models.AltNameFile, handler func(*models.AlternateName) error) error {
	return handle(p, archive.DumpFile(), true, handler)
}

func (p Parser) GetLanguages(handler func(*models.Language) error) error {
	return handle(p, LangCodes, false, handler)
}

func (p Parser) GetTimeZones(handler func(*models.TimeZone) error) error {
	return handle(p, TimeZones, false, handler)
}

func (p Parser) GetCountries(handler func(*models.Country) error) error {
	return handle(p, Countries, true, handler)
}

func (p Parser) GetFeatureCodes(file models.FeatureCodeFile, handler func(*models.FeatureCode) error) error {
	return handle(p, models.DumpFile(file), true, handler)
}

func (p Parser) GetHierarchy(handler func(*models.Hierarchy) error) error {
	return handle(p, Hierarchy, true, handler)
}

func (p Parser) GetShapes(handler func(*models.Shape) error) error {
	return handle(p, Shapes, false, handler)
}

func (p Parser) GetUserTags(handler func(*models.UserTag) error) error {
	return handle(p, UserTags, true, handler)
}

func (p Parser) GetAdminDivisions(handler func(*models.AdminDivision) error) error {
	return handle(p, AdminDivisions, true, handler)
}

func (p Parser) GetAdminSubdivisions(handler func(*models.AdminSubdivision) error) error {
	return handle(p, AdminSubDivisions, true, handler)
}

func (p Parser) GetAdminCodes5(handler func(*models.AdminCode5) error) error {
	return handle(p, AdminCode5, true, handler)
}

func (p Parser) GetAlternateNameDeletes(handler func(*models.AlternateNameDelete) error) error {
	return handle(p, AlternateNamesDeletes.WithLastDate(), true, handler)
}

func (p Parser) GetAlternateNameModifications(handler func(*models.AlternateNameModification) error) error {
	return handle(p, AlternateNamesModifications.WithLastDate(), true, handler)
}

func (p Parser) GetDeletes(handler func(*models.GeonameDelete) error) error {
	return handle(p, Deletes.WithLastDate(), true, handler)
}

func (p Parser) GetModifications(handler func(*models.Geoname) error) error {
	return handle(p, Modifications.WithLastDate(), true, handler)
}
