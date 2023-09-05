package main

import (
	"fmt"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/v3v3r3v/geonames"
	"github.com/v3v3r3v/geonames/models"
	"os"
	"runtime"
	"time"
)

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func main() {
	w := wow.New(os.Stdout, spin.Spinner{Frames: []string{"⚙️"}}, "  Parsing all cities with a population > 5000...")
	w.Persist()

	remoteFetcher := geonames.NewFetcher(geonames.FetcherConfig{
		RemoteUrl: geonames.DownloadGeonamesOrgUrl,
	})

	for _, file := range []models.DumpFile{
		geonames.Cities5000.DumpFile(),
		geonames.AlternateNames.DumpFile(),
	} {
		if fileExists(file.String()) {
			w.PersistWith(spin.Spinner{Frames: []string{"✅"}}, fmt.Sprintf(" Local dump file found: %s", file.String()))
		} else {
			w.Text("").Spinner(spin.Get(spin.Earth)).Start()
			w.Text(fmt.Sprintf("Fetching remote dump: %s", file.String()))
			err := remoteFetcher.DumpToFile(file)
			if err != nil {
				w.PersistWith(spin.Spinner{Frames: []string{"🔥"}}, fmt.Sprintf(" Save Dump %s Error: %s", file.String(), err.Error()))
				return
			}
			w.PersistWith(spin.Spinner{Frames: []string{"✅"}}, fmt.Sprintf(" Local dump file %s saved!", file.String()))
		}
	}

	localFetcher := geonames.NewFetcher(geonames.FetcherConfig{})

	parser := geonames.Parser{
		FetchSource: geonames.SourceFs,
		Fetcher:     localFetcher,
	}

	w.Text("").Spinner(spin.Get(spin.Earth)).Start()
	count := 0
	since := time.Now()
	var m runtime.MemStats
	var max uint64 = 0
	err := parser.GetGeonames(geonames.Cities5000, func(alternateName *models.Geoname) error {
		count++
		if count%10000 == 0 {
			w.Text(fmt.Sprintf("%d: %s", count, alternateName.Name))
		}

		runtime.ReadMemStats(&m)
		if max < m.Alloc {
			max = m.Alloc
		}
		return nil
	})
	if err != nil {
		w.PersistWith(spin.Spinner{Frames: []string{"🔥"}}, fmt.Sprintf(" Error: %s", err.Error()))
		return
	}
	duration := time.Since(since)

	w.PersistWith(spin.Spinner{Frames: []string{"✅"}}, fmt.Sprintf(" Done!"))
	w.PersistWith(spin.Spinner{Frames: []string{"⛩"}}, fmt.Sprintf("  Cities: %d", count))
	w.PersistWith(spin.Spinner{Frames: []string{"⏱"}}, fmt.Sprintf("  Duration: %d sec", duration/time.Second))
	w.PersistWith(spin.Spinner{Frames: []string{"💾️‍"}}, fmt.Sprintf(" Memory: %d Mb", max/1024/1024))
}
