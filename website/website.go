/*
Copyright (C) 2015, 2016 NanoDano <nanodano@devdungeon.com>

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/DevDungeon/WebGenome/core"
	"github.com/dustin/go-humanize"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	resultsPerPage int = 25
)

func getParentDomains(domain core.Domain, dbConn *mgo.Collection) []core.Domain {
	var parentDomains []core.Domain
	var tempDomain core.Domain
	var emptyId bson.ObjectId

	// End condition for recursion
	if domain.ParentDomain == emptyId {
		return parentDomains // Empty
	}

	getParentQuery := bson.M{"_id": domain.ParentDomain}
	dbConn.Find(getParentQuery).One(&tempDomain)

	// Return this domain as the last parent or recurse deeper
	if tempDomain.ParentDomain == emptyId {
		parentDomains = append(parentDomains, tempDomain)
	} else {
		parentDomains = append(parentDomains, tempDomain)
		parentDomains = append(parentDomains, getParentDomains(tempDomain, dbConn)...)
	}
	return parentDomains
}

func viewDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer := render.New(render.Options{
		Layout: "layout",
	})
	var domain core.Domain
	var parentDomains []core.Domain

	query := bson.M{"_id": bson.ObjectIdHex(ps.ByName("id"))}

	session, _ := mgo.Dial("localhost")
	defer session.Close()
	dbConn := session.DB("webgenome").C("domains")
	dbConn.Find(query).One(&domain)

	// Get all parents
	parentDomains = getParentDomains(domain, dbConn)

	vars := map[string]interface{}{
		"title":         "View Domain",
		"domain":        domain,
		"parentDomains": parentDomains,
	}

	renderer.HTML(w, http.StatusOK, "view_domain", vars)
}

//func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"", ""}}}}
//	renderDomainList(w, r, p, query, "")
//}

func drupal(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Drupal", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Drupal Sites")
}

func django(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^django_language", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Django Sites")
}

func zope(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^Zope/", "i"}}}}
	renderDomainListFromQuery(w, r, p, query, "Zope Sites")
}

func php(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^PHP/", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "PHP Sites")
}

func java(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^JSESSIONID=", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Java Sites")
}

func aspdotnet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^ASP.NET_SessionId=", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "ASP.NET Sites")
}

func python(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Python/", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Python Sites")
}

func ruby(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Ruby/", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Ruby Sites")
}

func apache(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^Apache", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Apache Sites")
}

func nginx(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^nginx", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Nginx Sites")
}

func iis(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^Microsoft-IIS", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "IIS Sites")
}

func tomcat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"tomcat", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Tomcat Sites")
}

func webrick(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^WEBrick", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "WEBrick Sites")
}

func lighttpd(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^lighttpd", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Lighttpd Sites")
}

func ibmhttpserver(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"^IBM_HTTP_Server", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "IBM HTTP Server")
}

func apusic(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Apusic", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Apusic Sites")
}

func enhydra(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Enhydra", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Enhydra Sites")
}

func jetty(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Jetty", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Jetty Sites")
}

func unix(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"(Unix)", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Unix Sites")
}

func linux(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Linux", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Linux Sites")
}

func debian(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Debian", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Debian sites")
}

func fedora(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Fedora", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Fedora Sites")
}

func redhat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Red Hat", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Red Hat Sites")
}

func centos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"CentOS", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "CentOS Sites")
}

func ubuntu(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Ubuntu", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Ubuntu Sites")
}

func freebsd(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"FreeBSD", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "FreeBSD Sites")
}

func win32(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Win32", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Win32 Sites")
}

func win64(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Win64", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Win64 Sites")
}

func darwin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Darwin", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Darwin Sites")
}

func phusionpassenger(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Phusion_Passenger", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Phusion Passenger Sites")
}

func openssl(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"OpenSSL", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "OpenSSL Sites")
}

func webdav(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"DAV", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "WebDAV Sites")
}

func communique(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"Communique", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "Communique Sites")
}

func bigipserver(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"headers": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{"BIGipServer", ""}}}}
	renderDomainListFromQuery(w, r, p, query, "BIGipServer Sites")
}

func gov(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := bson.M{"name": bson.RegEx{".gov", ""}}
	renderDomainListFromQuery(w, r, p, query, "Government Sites")
}

func renderDomainListFromQuery(w http.ResponseWriter, r *http.Request, _ httprouter.Params, query bson.M, title string) {

	var (
		page         string
		pageNumber   int
		previousPage string
		nextPage     string
	)
	page = r.URL.Query().Get("page")
	if page == "" {
		pageNumber = 1
	} else {
		pageNumber, _ = strconv.Atoi(page)
	}

	var domains []core.Domain
	session, _ := mgo.Dial("localhost")
	defer session.Close()
	dbConn := session.DB("webgenome").C("domains")
	dbConn.Find(query).Limit(resultsPerPage).Skip((pageNumber - 1) * resultsPerPage).All(&domains)

	if pageNumber <= 1 {
		previousPage = ""
	} else {
		previousPage = r.URL.Path + "?page=" + strconv.Itoa(pageNumber-1)
	}
	if len(domains) < resultsPerPage {
		nextPage = ""
	} else {
		nextPage = r.URL.Path + "?page=" + strconv.Itoa(pageNumber+1)
	}

	vars := map[string]interface{}{
		"title":        title,
		"domains":      domains,
		"pageNumber":   pageNumber,
		"previousPage": previousPage,
		"nextPage":     nextPage,
	}
	renderer := render.New(render.Options{
		Layout: "layout",
	})
	renderer.HTML(w, http.StatusOK, "domain_listing", vars)
}

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session, _ := mgo.Dial("localhost")
	defer session.Close()
	dbConn := session.DB("webgenome").C("domains")
	totalDomains, err := dbConn.Count()
	if err != nil {
		fmt.Println("Error getting total domain count." + err.Error())
	}

	vars := map[string]interface{}{
		"title":        "Home",
		"totalDomains": humanize.Comma(int64(totalDomains)),
	}

	renderer := render.New(render.Options{
		Layout: "layout",
	})
	renderer.HTML(w, http.StatusOK, "index", vars)
}

func random(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session, _ := mgo.Dial("localhost")
	defer session.Close()
	dbConn := session.DB("webgenome").C("domains")
	totalDomains, err := dbConn.Count()
	if err != nil {
		fmt.Println("Error getting total domain count." + err.Error())
	}

	query := bson.M{"headers": bson.M{"$exists": true}}
	var domain core.Domain
	numDomainsToSkip := rand.Intn(totalDomains) % 1000000 // Hard limit of 5 mill for speed
	err = dbConn.Find(query).Limit(1).Skip(numDomainsToSkip).One(&domain)
	if err != nil {
		fmt.Println("Random lookup went too far. Redirecting back to random.")
		http.Redirect(w, r, "/random", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/domain/"+domain.Id.Hex(), http.StatusFound)
	return
}

func search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		pageNumber int
		nextPage   string
		domains    []core.Domain
	)
	session, _ := mgo.Dial("localhost")
	defer session.Close()
	dbConn := session.DB("webgenome").C("domains")

	domainKeyword := r.PostFormValue("domain-keyword")
	fmt.Println("Search query: " + domainKeyword)
	query := bson.M{"name": bson.RegEx{domainKeyword, "i"}}

	err := dbConn.Find(query).All(&domains)

	if err != nil {
		fmt.Println("Error with search.")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	pageNumber = 1

	dbConn.Find(query).Limit(resultsPerPage).All(&domains)

	if len(domains) < resultsPerPage {
		nextPage = ""
	} else {
		nextPage = "#"
	}

	vars := map[string]interface{}{
		"title":      "Search Results",
		"domains":    domains,
		"pageNumber": pageNumber,
		"nextPage":   nextPage,
	}
	renderer := render.New(render.Options{
		Layout: "layout",
	})
	renderer.HTML(w, http.StatusOK, "search_results", vars)

}

func loggerMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("[*] " + time.Now().Format(time.RFC850) + " - " + r.Header.Get("X-Real-IP"))
	// do some IP logging
	next(rw, r)
	// Do some stuff after
}

func premium(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	vars := map[string]interface{}{}

	renderer := render.New(render.Options{
		Layout: "layout",
	})
	renderer.HTML(w, http.StatusOK, "premium", vars)
}

func main() {
	staticFilesDir := "./static/"

	// Routing
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/", search)
	router.GET("/domain/:id", viewDomain)
	router.GET("/random", random)
	router.GET("/premium", premium)

	router.GET("/gov", gov)
	router.GET("/drupal", drupal)
	router.GET("/django", django)
	router.GET("/zope", zope)
	router.GET("/php", php)
	router.GET("/java", java)
	router.GET("/aspdotnet", aspdotnet)
	router.GET("/python", python)
	router.GET("/ruby", ruby)
	router.GET("/apache", apache)
	router.GET("/nginx", nginx)
	router.GET("/iis", iis)
	router.GET("/tomcat", tomcat)
	router.GET("/webrick", webrick)
	router.GET("/lighttpd", lighttpd)
	router.GET("/ibmhttpserver", ibmhttpserver)
	router.GET("/apusic", apusic)
	router.GET("/enhydra", enhydra)
	router.GET("/jetty", jetty)
	router.GET("/unix", unix)
	router.GET("/linux", linux)
	router.GET("/debian", debian)
	router.GET("/fedora", fedora)
	router.GET("/redhat", redhat)
	router.GET("/centos", centos)
	router.GET("/ubuntu", ubuntu)
	router.GET("/freebsd", freebsd)
	router.GET("/win32", win32)
	router.GET("/win64", win64)
	router.GET("/darwin", darwin)
	router.GET("/phusionpassenger", phusionpassenger)
	router.GET("/openssl", openssl)
	router.GET("/webdav", webdav)
	router.GET("/communique", communique)
	router.GET("/bigipserver", bigipserver)

	// Unchecked sites
	// Checked sites
	// Skipped sites
	// Total domains

	// Middleware
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewStatic(http.Dir(staticFilesDir)))
	n.Use(negroni.HandlerFunc(loggerMiddleware))

	n.UseHandler(router)
	n.Run(":3000")

}
