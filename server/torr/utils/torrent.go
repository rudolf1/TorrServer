package utils

import (
	"encoding/base32"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"server/config"
	"server/settings"

	"golang.org/x/time/rate"
)

var defTrackers = []string{
	"http://retracker.local",
	"http://bt4.t-ru.org/ann?magnet",
	"http://retracker.mgts.by:80/announce",
	"http://tracker.city9x.com:2710/announce",
	"http://tracker.electro-torrent.pl:80/announce",
	"http://tracker.internetwarriors.net:1337/announce",
	"http://tracker2.itzmx.com:6961/announce",
	"udp://opentor.org:2710",
	"udp://public.popcorn-tracker.org:6969/announce",
	"udp://tracker.opentrackr.org:1337/announce",
	"http://bt.svao-ix.ru/announce",
	"udp://explodie.org:6969/announce",
	"wss://tracker.btorrent.xyz",
	"wss://tracker.openwebtorrent.com",
}

var loadedTrackers []string

func GetTrackerFromFile() []string {
	buf, err := config.ReadConfigParser("Trackers")
	if err == nil {
		return buf
	} else {
		return nil
	}
}

func GetDefTrackers() []string {
	loadNewTracker()
	if len(loadedTrackers) == 0 {
		return defTrackers
	}
	return loadedTrackers
}

func RemoveDuplicates[T string | int](tSlice []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range tSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RemoveLine(a []string, i int) {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = ""
	a = a[:len(a)-1]
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

var defaultUrl = []string{
	"https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_best_ip.txt",
	//	"https://newtrackon.com/api/stable",
	//	"https://raw.githubusercontent.com/XIU2/TrackersListCollection/master/best.txt",
}

func loadNewTracker() {
	if len(loadedTrackers) > 0 {
		return
	}
	got, get := config.ReadConfigParser("Default_url")
	if len(got) > 0 && get == nil {
		defaultUrl = nil
		defaultUrl = got
	}
	for _, a := range defaultUrl {
		var resp, err = http.Get(a)
		if err == nil {
			buf, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				arr := strings.Split(string(buf), "\n")
				var ret []string
				for _, s := range arr {
					s = strings.TrimSpace(s)
					if len(s) > 0 {
						ret = append(ret, s)
					}
				}
				loadedTrackers = append(loadedTrackers, ret...)
			}
		}
	}
	loadedTrackers = append(loadedTrackers, defTrackers...)
	TrackersDel, back := config.ReadConfigParser("Blacklist_tracker")
	if back == nil {
		if len(TrackersDel) > 0 {
			for _, a := range TrackersDel {
				if Contains(loadedTrackers, a) {
					i := Find(loadedTrackers, a)
					RemoveLine(loadedTrackers, i)
				}
			}
		}
	}
	loadedTrackers = RemoveDuplicates(loadedTrackers)
	path := filepath.Join(settings.Path, "trackers.tmp")
	file, err := os.Create(path)
	if err == nil {
		defer file.Close()
		for _, c := range loadedTrackers {
			_, err := file.WriteString(c + "\n")
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println("Trackers file done")
	}
}

func PeerIDRandom(peer string) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return peer + base32.StdEncoding.EncodeToString(randomBytes)[:20-len(peer)]
}

func Limit(i int) *rate.Limiter {
	l := rate.NewLimiter(rate.Inf, 0)
	if i > 0 {
		b := i
		if b < 16*1024 {
			b = 16 * 1024
		}
		l = rate.NewLimiter(rate.Limit(i), b)
	}
	return l
}
