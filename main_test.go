package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testURLsPair struct {
	Owner, Repo string
	URLs        map[string]string
}

func TestBuildURLs(t *testing.T) {
	tps := []testURLsPair{
		{"spectolabs", "hoverfly", map[string]string{
			"repo":   "https://api.github.com/repos/spectolabs/hoverfly",
			"views":  "https://api.github.com/repos/spectolabs/hoverfly/traffic/views",
			"clones": "https://api.github.com/repos/spectolabs/hoverfly/traffic/clones",
			"rels":   "https://api.github.com/repos/spectolabs/hoverfly/releases",
		},
		},
		{"spectolabs", "hoverfly-java", map[string]string{
			"repo":   "https://api.github.com/repos/spectolabs/hoverfly-java",
			"views":  "https://api.github.com/repos/spectolabs/hoverfly-java/traffic/views",
			"clones": "https://api.github.com/repos/spectolabs/hoverfly-java/traffic/clones",
			"rels":   "https://api.github.com/repos/spectolabs/hoverfly-java/releases",
		},
		},
		{"golang", "go", map[string]string{
			"repo":   "https://api.github.com/repos/golang/go",
			"views":  "https://api.github.com/repos/golang/go/traffic/views",
			"clones": "https://api.github.com/repos/golang/go/traffic/clones",
			"rels":   "https://api.github.com/repos/golang/go/releases",
		},
		},
	}
	for _, tp := range tps {
		u := buildURLs(tp.Owner, tp.Repo)
		assert.Equal(t, tp.URLs, u)
	}
}

func TestRequestBytes(t *testing.T) {
	bytes, _ := request("http://ip.jsontest.com")
	assert.NotEqual(t, len(bytes), 0, "bytes length is 0")
}

func TestRequestNotFound(t *testing.T) {
	u := "http://hoverfly.io/blah"
	_, err := request(u)
	assert.Error(t, err, "an error was expected")
}

type testRepoDataPair struct {
	URL string
	testRepoData
}

type testRepoData struct {
	Stars, Forks, Issues int
}

func TestGetRepoData(t *testing.T) {
	tps := []testRepoDataPair{
		testRepoDataPair{
			"https://api.github.com/repos/spectolabs/hoverfly",
			testRepoData{663, 50, 20},
		},
		testRepoDataPair{
			"https://api.github.com/repos/spectolabs/hoverfly-java",
			testRepoData{22, 5, 10},
		},
		testRepoDataPair{
			"https://api.github.com/repos/golang/go",
			testRepoData{24132, 3228, 2595},
		},
	}
	for _, tp := range tps {
		rd, _ := getRepoData(tp.URL)
		assert.Equal(t, tp.testRepoData.Stars, rd.Stars)
		assert.Equal(t, tp.testRepoData.Forks, rd.Forks)
		assert.Equal(t, tp.testRepoData.Issues, rd.Issues)
	}
}

type testViewsDataPair struct {
	URL string
	viewsData
}

type viewsData struct {
	LatestCount, LatestUniques int
}

func TestGetViewsData(t *testing.T) {
	d := time.Date(2017, 1, 25, 0, 0, 0, 0, time.UTC)
	tps := []testViewsDataPair{
		testViewsDataPair{
			"https://api.github.com/repos/spectolabs/hoverfly/traffic/views",
			viewsData{46, 34},
		},
		testViewsDataPair{
			"https://api.github.com/repos/spectolabs/hoverfly-java/traffic/views",
			viewsData{3, 3},
		},
		testViewsDataPair{
			"https://api.github.com/repos/spectolabs/hoverpy/traffic/views",
			viewsData{1, 1},
		},
	}
	for _, tp := range tps {
		vd, _ := getViewsData(tp.URL, d)
		assert.Equal(t, tp.viewsData.LatestCount, vd.LatestCount)
		assert.Equal(t, tp.viewsData.LatestUniques, vd.LatestUniques)
	}
}

type testClonesDataPair struct {
	URL string
	clonesData
}

type clonesData struct {
	LatestCount, LatestUniques int
}

func TestGetClonesData(t *testing.T) {
	d := time.Date(2017, 1, 25, 0, 0, 0, 0, time.UTC)
	tps := []testClonesDataPair{
		testClonesDataPair{
			"https://api.github.com/repos/spectolabs/hoverfly/traffic/clones",
			clonesData{1, 1},
		},
		testClonesDataPair{
			"https://api.github.com/repos/spectolabs/hoverfly-java/traffic/clones",
			clonesData{0, 0},
		},
		testClonesDataPair{
			"https://api.github.com/repos/spectolabs/hoverpy/traffic/clones",
			clonesData{0, 0},
		},
	}
	for _, tp := range tps {
		cd, _ := getClonesData(tp.URL, d)
		assert.Equal(t, tp.clonesData.LatestCount, cd.LatestCount)
		assert.Equal(t, tp.clonesData.LatestUniques, cd.LatestUniques)
	}
}

type testDownloadDataPair struct {
	URL  string
	Dlds int
}

func TestGetDownloadCount(t *testing.T) {
	tps := []testDownloadDataPair{
		testDownloadDataPair{
			"https://api.github.com/repos/spectolabs/hoverfly/releases", 2845},
		testDownloadDataPair{
			"https://api.github.com/repos/spectolabs/hoverfly-java/releases", 0},
		testDownloadDataPair{
			"https://api.github.com/repos/golang/go/releases", 0},
	}
	for _, tp := range tps {
		dc, _ := getDownloadCount(tp.URL)
		assert.Equal(t, tp.Dlds, dc)
	}
}

type testDataPair struct {
	Owner, Repo string
	Day         time.Time
	testData
}

type testData struct {
	Owner, Repo                                                               string
	Date                                                                      time.Time
	Stars, Forks, Issues, Views, UniqueViews, Clones, UniqueClones, Downloads int
}

func TestGetData(t *testing.T) {
	d := time.Date(2017, 1, 27, 0, 0, 0, 0, time.UTC)
	tps := []testDataPair{
		testDataPair{
			"spectolabs", "hoverfly", d,
			testData{"spectolabs", "hoverfly", d, 663, 50, 20, 27, 20, 7, 4, 2845},
		},
		testDataPair{
			"spectolabs", "hoverfly-java", d,
			testData{"spectolabs", "hoverfly-java", d, 22, 5, 10, 9, 6, 0, 0, 0},
		},
		testDataPair{
			"spectolabs", "hoverpy", d,
			testData{"spectolabs", "hoverpy", d, 66, 2, 4, 0, 0, 1, 1, 0},
		},
	}
	for _, tp := range tps {
		ds, _ := getData(tp.Owner, tp.Repo, tp.Day)
		assert.Equal(t, tp.Owner, ds.Owner)
		assert.Equal(t, tp.Repo, ds.Repo)
		assert.Equal(t, tp.Date, ds.Date)
		assert.Equal(t, tp.Stars, ds.Stars)
		assert.Equal(t, tp.Forks, ds.Forks)
		assert.Equal(t, tp.Issues, ds.Issues)
		assert.Equal(t, tp.Views, ds.Views)
		assert.Equal(t, tp.UniqueViews, ds.UniqueViews)
		assert.Equal(t, tp.Clones, ds.Clones)
		assert.Equal(t, tp.UniqueClones, ds.UniqueClones)
		assert.Equal(t, tp.Downloads, ds.Downloads)
	}
}
