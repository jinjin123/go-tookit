package lib

import (
	"strings"
)

type Pager struct {
	Page     int
	Totalnum int
	Pagesize int
	urlpath  string
	urlquery string
	nopath   bool
}

/*  page  parse https://so.csdn.net/so/search/all?q=query&t=all&page=50
?q=query&t=all&page=50 */
func NewPager(page, totalnum, pagesize int, url string, nopath ...bool) *Pager {
	p := new(Pager)
	p.Page = page
	p.Totalnum = totalnum
	p.Pagesize = pagesize

	arr := strings.Split(url, "?")
	p.urlpath = arr[0]
	if len(arr) > 1 {
		p.urlquery = "?" + arr[1]
	} else {
		p.urlquery = ""
	}

	if len(nopath) > 0 {
		p.nopath = nopath[0]
	} else {
		p.nopath = false
	}

	return p
}
