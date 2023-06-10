package poster_tmdb

import (
	"github.com/StalkR/imdb"
	"net/http"
	"regexp"
	"server/log"
	"strings"
	"time"
)

var tor_q = []string{
	"Blu-ray",
	"BDRemux",
	"BDRip",
	"HDRip",
	"WEB-DL",
	"WEBDL",
	"XviD",
	"WEB-DLRip",
	"WEBRip",
	"HDTV",
	"HDTVRip",
	"DVD",
	"DVD9",
	"DVD5",
	"DVDRip",
	"DVDScr",
	"DVDRemux",
	"DVB",
	"SATRip",
	"IPTVRip",
	"TVRip",
	"VHSRip",
	"TS",
	"CAMRip",
	"2160p",
	"1080p",
	"1080i",
	"720p",
	"576p",
	"576i",
	"480p",
	"400p",
}

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"

type customTransport struct {
	http.RoundTripper
}

func (e *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	defer time.Sleep(time.Second)         // don't go too fast or risk being blocked by awswaf
	r.Header.Set("Accept-Language", "en") // avoid IP-based language detection
	r.Header.Set("User-Agent", userAgent)
	return e.RoundTripper.RoundTrip(r)
}

func GetPoster(name string) string {
	name = strings.ReplaceAll(name, ".", " ")
	name = strings.ReplaceAll(name, "|", "")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	name = strings.ReplaceAll(name, " RUS ", " ")
	name = strings.ReplaceAll(name, " ENG ", " ")
	nameMass := strings.Split(name, " ")
	gp, err := regexp.Compile("[0-9][0-9][0-9][0-9]")
	if err != nil {
		log.TLogln("Error compile regex %v", err)
	}
	for i, word := range nameMass {
		for _, word2 := range tor_q {
			if strings.EqualFold(word, word2) {
				for l := i; l < len(nameMass); l++ {
					nameMass[l] = ""
				}
				nameMass = nameMass[:i]
			}
		}
		if len(gp.FindString(word)) > 0 {
			for m := i + 1; m < len(nameMass); m++ {
				nameMass[m] = ""
			}
			nameMass = nameMass[:i+1]
		}
	}
	nameMassNew := strings.Join(nameMass, " ")
	if len(nameMassNew) > 0 {
		nameMassNew = strings.Trim(nameMassNew, " ")
	}
	gp, err = regexp.Compile("\\[[a-zA-Zа-яА-Я0-9-[:space:]+.,].+\\]")
	if err != nil {
		log.TLogln("Error compile regex %v", err)
	}
	nameMassNew = strings.ReplaceAll(nameMassNew, gp.FindString(nameMassNew), "")
	gp, err = regexp.Compile("[Ss][0-9][0-9]")
	if err != nil {
		log.TLogln("Error compile regex %v", err)
	}
	nameMassNew = strings.ReplaceAll(nameMassNew, gp.FindString(nameMassNew), "")
	gp, err = regexp.Compile("[Ee][0-9][0-9]")
	if err != nil {
		log.TLogln("Error compile regex %v", err)
	}
	nameMassNew = strings.ReplaceAll(nameMassNew, gp.FindString(nameMassNew), "")
	if strings.Contains(nameMassNew, "/") {
		nameMass = strings.Split(nameMassNew, "/")
		gp, err = regexp.Compile("[A-Za-z]+")
		if err != nil {
			log.TLogln("Error compile regex %v", err)
		}
		for _, word := range nameMass {
			out := gp.FindString(word)
			if len(out) > 2 {
				nameMassNew = word
				break
			}
		}
		if len(nameMassNew) > 0 {
			nameMassNew = strings.Trim(nameMassNew, " ")
		} else {
			for _, word2 := range nameMass {
				nameMassNew = word2
				break
			}
			if len(nameMassNew) > 0 {
				nameMassNew = strings.Trim(nameMassNew, " ")
			}
		}
	}
	log.TLogln(nameMassNew)
	client := &http.Client{
		Transport: &customTransport{http.DefaultTransport},
	}
	results, err := imdb.SearchTitle(client, nameMassNew)
	var imdb_id string
	if err == nil {
		log.TLogln("Result:", results[0].ID, results[0].Name)
		imdb_id = results[0].ID
	} else {
		log.TLogln(err)
		return ""
	}
	result, err2 := imdb.NewTitle(client, imdb_id)
	if err2 == nil {
		log.TLogln("Poster:", result.Poster.ContentURL)
		return result.Poster.ContentURL
	} else {
		log.TLogln(err2)
		return ""
	}
}
