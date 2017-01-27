package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	ghAPI   string = "https://api.github.com/repos/"
	vwsPath string = "/traffic/views"
	cnsPath string = "/traffic/clones"
	rlsPath string = "/releases"
)

var (
	user      = os.Getenv("GH_USER")
	token     = os.Getenv("GH_TOKEN")
	today     = time.Now().UTC().Truncate(24 * time.Hour)
	yesterday = today.AddDate(0, 0, -1)
)

func main() {
	d, err := getData("spectolabs", "hoverfly-java", yesterday)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)
}

type Repo struct {
	Owner, Repo                                                               string
	Date                                                                      time.Time
	Stars, Forks, Issues, Views, UniqueViews, Clones, UniqueClones, Downloads int
}

func buildURLs(owner, repo string) map[string]string {
	u := make(map[string]string)
	u["repo"] = ghAPI + owner + "/" + repo
	u["views"] = u["repo"] + vwsPath
	u["clones"] = u["repo"] + cnsPath
	u["rels"] = u["repo"] + rlsPath
	return u
}

func getData(owner, repo string, day time.Time) (*Repo, error) {
	urls := buildURLs(owner, repo)
	rd, err := getRepoData(urls["repo"])
	if err != nil {
		return nil, err
	}
	vc, err := getViewsData(urls["views"], day)
	if err != nil {
		return nil, err
	}
	cc, err := getClonesData(urls["clones"], day)
	if err != nil {
		return nil, err
	}
	dc, err := getDownloadCount(urls["rels"])
	if err != nil {
		return nil, err
	}
	r := new(Repo)
	r.Owner = owner
	r.Repo = repo
	r.Date = today
	r.Stars = rd.Stars
	r.Forks = rd.Forks
	r.Issues = rd.Issues
	r.Views = vc.LatestCount
	r.UniqueViews = vc.LatestUniques
	r.Clones = cc.LatestCount
	r.UniqueClones = cc.LatestUniques
	r.Downloads = dc

	return r, err

}

func request(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, token)
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		err := fmt.Errorf("%v returned %v", url, res.StatusCode)
		return nil, err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

type RepoData struct {
	Stars  int `json:"stargazers_count"`
	Forks  int `json:"forks_count"`
	Issues int `json:"open_issues_count"`
}

func getRepoData(url string) (*RepoData, error) {
	rd := new(RepoData)
	b, err := request(url)
	json.Unmarshal(b, &rd)
	return rd, err
}

type ViewsData struct {
	Views []struct {
		Timestamp time.Time `json:"timestamp"`
		Count     int       `json:"count"`
		Uniques   int       `json:"uniques"`
	} `json:"views"`
	LatestCount   int
	LatestUniques int
}

func getViewsData(url string, day time.Time) (*ViewsData, error) {
	vs := new(ViewsData)
	b, err := request(url)
	json.Unmarshal(b, &vs)
	for _, d := range vs.Views {
		if d.Timestamp == day {
			vs.LatestCount = d.Count
			vs.LatestUniques = d.Uniques
		}
	}
	return vs, err
}

type ClonesData struct {
	Clones []struct {
		Timestamp time.Time `json:"timestamp"`
		Count     int       `json:"count"`
		Uniques   int       `json:"uniques"`
	} `json:"clones"`
	LatestCount   int
	LatestUniques int
}

func getClonesData(url string, day time.Time) (*ClonesData, error) {
	cs := new(ClonesData)
	b, err := request(url)
	json.Unmarshal(b, &cs)
	for _, d := range cs.Clones {
		if d.Timestamp == day {
			cs.LatestCount = d.Count
			cs.LatestUniques = d.Uniques
		}
	}
	return cs, err
}

type ReleasesData []struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		Name          string `json:"name"`
		DownloadCount int    `json:"download_count"`
	} `json:"assets"`
}

func getDownloadCount(url string) (int, error) {
	var dls int
	var rs ReleasesData
	b, err := request(url)
	json.Unmarshal(b, &rs)
	for _, r := range rs {
		for _, a := range r.Assets {
			dls += a.DownloadCount
		}
	}
	return dls, err
}
