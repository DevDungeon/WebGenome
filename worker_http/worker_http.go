package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/DevDungeon/WebGenome/core"

	"github.com/PuerkitoBio/goquery"
	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	verbose bool
)

// Ignores verbosity option
func logError(message string) {
	color.Set(color.FgRed)
	log.Println("[!] " + message)
	color.Unset()
}

// Ignores verbosity option
func logGreen(message string) {
	color.Set(color.FgGreen)
	log.Println("[+] " + message)
	color.Unset()
}

func logWarning(message string) {
	if verbose {
		color.Set(color.FgYellow)
		log.Println("[-] " + message)
		color.Unset()
	}
}

func logInfo(message string) {
	if verbose {
		color.Set(color.FgCyan)
		log.Println("[*] " + message)
		color.Unset()
	}
}

func check(err error) {
	if err != nil {
		logError("Fatal error: " + err.Error())
		os.Exit(1)
	}
}

func appendIfNotExists(strings []string, newString string) []string {
	exists := false
	for _, existingString := range strings {
		if existingString == newString {
			exists = true
		}
	}
	if !exists {
		strings = append(strings, newString)
	}
	return strings
}

// Extract URLs from HTML document
func getUniqueDomainsFromResponse(response *http.Response) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	// Get all a hrefs
	var domains []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			cleanDomain, found := getDomainFromHref(href)
			if found {
				domains = appendIfNotExists(domains, cleanDomain)
			}
		}
	})
	return domains, nil
}

// Clean up an href and find the root domain if available
// handle things like: protocol prefixes, mailto links, trailing slashes
// trailing hash
func getDomainFromHref(href string) (domain string, found bool) {
	// If it is too short it can't be a domain and is probably just a /
	if len(href) < 3 {
		return "", false
	}

	// If it starts right off with a slash, it is a relative URL and no domain
	if href[0] == '/' {
		return "", false
	}

	// If there are double slashes, strip them and everything before
	pos := strings.Index(href, "//")
	if pos > -1 {
		href = href[pos+2:]
	} else {
		return "", false // No protocol, treat as relative and skip
	}

	// If it is a mailto, then strip the mailto,
	pos = strings.Index(href, "mailto:")
	if pos > -1 {
		href = href[pos+7:]
	}

	// Remove any spaces and anything after
	pos = strings.Index(href, " ")
	if pos > -1 {
		href = href[:pos]
	}

	// Remove any username@ portion
	pos = strings.Index(href, "@")
	if pos > -1 {
		href = href[:pos]
	}

	// Remove any ? and after
	pos = strings.Index(href, "?")
	if pos > -1 {
		href = href[:pos]
	}

	// Remove any & and after
	pos = strings.Index(href, "&")
	if pos > -1 {
		href = href[:pos]
	}

	// Remove any % and after
	pos = strings.Index(href, "%")
	if pos > -1 {
		href = href[:pos]
	}

	// If there are still any more slashes they are at the end
	// so trim them and anything after
	pos = strings.Index(href, "/")
	if pos > -1 {
		href = href[:pos]
	}

	// Strip # and everything after
	pos = strings.Index(href, "#")
	if pos > -1 {
		href = href[:pos]
	}

	if !validateDomain(href) {
		return "", false
	}

	// if no good domain name obtained, return false
	return strings.ToLower(href), true
}

func validateDomain(domain string) bool {
	if len(domain) < 3 {
		return false
	}

	splitDomain := strings.Split(domain, ".")
	if len(splitDomain) < 2 { // There is not two parts to the domain
		return false
	}

	if len(splitDomain[0]) < 1 {
		return false
	}
	if len(splitDomain[1]) < 2 {
		return false
	}

	return true
}

func getDotCount(text string) int {
	// How many periods are in the text
	var dotCount int = 0
	for i := 0; i < len(text); i++ {
		if string(text[i]) == "." {
			dotCount++
		}
	}
	return dotCount
}

func areFirstFourLettersWwwDot(text string) bool {
	if len(text) < 4 {
		return false
	}
	return text[0:4] == "www."
}

func processDomain(domain core.Domain, httpTimeout time.Duration, doneChannel chan bool, dbConn *mgo.Collection) {

	var (
		err       error
		userAgent string = "WebGeno.me Research Project - Report issues to support@webgeno.me"
	)
	domain.LastChecked = time.Now()

	// Ignore some subdomains. These are like black holes with almost infinite subdomains
	// TODO Move this to a config file
	ignoredDomains := []string{
		".blogspot.com",
		".tumblr.com",
		".booked.net",
		".deviantart.com",
		".zxdyw.com",
		".fang.com",
		".8671.net",
	}
	for _, ignoredDomain := range ignoredDomains {
		pos := strings.Index(domain.Name, ignoredDomain)
		if pos > -1 {
			logInfo("Skipping ignored subdomain: " + domain.Name)
			domain.Skipped = true
			err = dbConn.Update(
				bson.M{"_id": domain.Id},
				domain,
			)
			check(err)
			doneChannel <- true
			return
		}
	}

	transport := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{
		Transport: transport,
		Timeout:   httpTimeout,
	}
	request, err := http.NewRequest("GET", "http://"+domain.Name+"/", nil)
	if err != nil {
		logError("Problem creating HTTP GET request for " + domain.Name + ". Setting skipped.")
		domain.Skipped = true
		err = dbConn.Update(
			bson.M{"_id": domain.Id},
			domain,
		)
		check(err)
		doneChannel <- true
		return
	}

	request.Header.Set("User-Agent", userAgent)
	request.Close = true
	request.Header.Set("Connection", "close") // Double check the connection is closed
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		domain.Skipped = true
		logWarning("Problem with " + domain.Name + ". Setting skipped. " + err.Error())
		if strings.Contains(err.Error(), "too many open files") {
			logError("Detecting too many files open error. Waiting 30 seconds.")
			time.Sleep(30 * time.Second)
			logInfo("Thread done waiting for 30 seconds.")
		}
		err = dbConn.Update(
			bson.M{"_id": domain.Id},
			domain,
		)
		check(err)
		doneChannel <- true
		return
	}

	// Pull out the headers from the HTTP response
	for key, value := range response.Header {
		if key == "Date" { // Ignore the Date header
			continue
		}
		header := core.Header{key, value[0]}
		domain.Headers = append(domain.Headers, header)
	}

	// Update domain
	err = dbConn.Update(
		bson.M{"_id": domain.Id},
		domain,
	)
	check(err)
	logInfo("Updated domain info: " + domain.Name)

	// Parse body for new domains
	domainsInDocument, err := getUniqueDomainsFromResponse(response)
	if err != nil {
		logInfo("Error reading response from: " + domain.Name)
	}
	logInfo("Domains found in " + domain.Name + ": " + strings.Join(domainsInDocument, ","))

	// Add new domains to database if they don't already exist
	for _, foundDomainName := range domainsInDocument {

		// See if domain exists
		// find count? if count > 0 or exists, move on
		var existingDomain core.Domain
		err = dbConn.Find(bson.M{"name": foundDomainName}).One(&existingDomain)
		if err == mgo.ErrNotFound {
			//logInfo("Domain not found. Adding: " + foundDomainName) // Just too verbose
			newDomain := &core.Domain{Name: foundDomainName, ParentDomain: domain.Id}
			err := dbConn.Insert(newDomain)
			if err != nil {
				logError("Error inserting!")
			}
		} else if err != nil {
			logError("Error when looking for existing domain: " + foundDomainName + ". " + err.Error())
		}
	}

	doneChannel <- true
	return
}

// Find unchecked domains (no headers)
func getDomainsToCheck(limit int, dbConn *mgo.Collection) ([]core.Domain, error) {
	var uncheckedDomains []core.Domain
	err := dbConn.Find(
		bson.M{"headers": bson.M{"$exists": 0}, "skipped": bson.M{"$exists": 0}},
	).Limit(limit).All(&uncheckedDomains)
	return uncheckedDomains, err

}

func main() {
	usage := `worker_http - Web Genome HTTP Worker.

Usage:
  worker_http --host=<host> --database=<dbname> --collection=<collectionname> --max-threads=<maxthreads> --http-timeout=<seconds> --batch-size=<batchsize> [--verbose]
  worker_http -h | --help
  worker_http --version

Options:
  -h --help                   Show this screen.
  --version                   Show version.
  --host=<host>               MongoDB host
  --database=<database>       MongoDB database name.
  --collection=<collection>   MongoDB collection name.
  --max-threads=<maxthreads>  Maximum number of simultaneous threads.
  --http-timeout=<seconds>    How long before HTTP requests timeout in seconds.
  --batch-size=<batchsize>    How many unchecked domains to pull and run per loop
  --verbose                   Increase output verbosity.`

	arguments, err := docopt.Parse(usage, nil, true, "Web Genome Worker", false)
	if err != nil {
		logError("Error parsing command line arguments. " + err.Error())
		os.Exit(1)
	}
	verbose = arguments["--verbose"].(bool) // Set global var for logging
	batchSize, err := strconv.Atoi(arguments["--batch-size"].(string))
	check(err)

	logGreen("====== Options ======")
	logGreen("Host:         " + arguments["--host"].(string))
	logGreen("Database:     " + arguments["--database"].(string))
	logGreen("Collection:   " + arguments["--collection"].(string))
	logGreen("Max threads:  " + arguments["--max-threads"].(string))
	logGreen("HTTP timeout: " + arguments["--http-timeout"].(string) + " seconds")
	logGreen("Batch size:   " + strconv.Itoa(batchSize))
	logGreen("Verbose Mode: " + strconv.FormatBool(verbose))
	logGreen("=====================")

	runtime.GOMAXPROCS(runtime.NumCPU())
	session, err := mgo.Dial(arguments["--host"].(string))
	check(err)
	defer session.Close()

	dbConn := session.DB(arguments["--database"].(string)).C(arguments["--collection"].(string))
	timeout, err := strconv.Atoi(arguments["--http-timeout"].(string))
	check(err)
	httpTimeout := time.Duration(time.Duration(timeout) * time.Second)
	maxThreads, err := strconv.Atoi(arguments["--max-threads"].(string))
	check(err)

	logGreen("Establishing connection with database.")
	logGreen("Database connection created.")

	totalRuns := 0
	doneChannel := make(chan bool)
	numThreads := 0
	for true {
		startTime := time.Now()
		uncheckedDomains, err := getDomainsToCheck(batchSize, dbConn)
		check(err)
		if len(uncheckedDomains) == 0 {
			logError("No domains found to check. Exiting.")
			os.Exit(0)
		}

		for x := 0; x < len(uncheckedDomains); x += 1 {
			numThreads += 1
			totalRuns += 1

			logGreen("Checking " + uncheckedDomains[x].Name)
			go processDomain(uncheckedDomains[x], httpTimeout, doneChannel, dbConn)

			// Wait until a done signal before next if max threads reached
			if numThreads >= maxThreads {
				<-doneChannel
				numThreads -= 1
			}
		}

		// Wait for all threads before repeating and fetching a new batch
		for numThreads > 0 {
			<-doneChannel
			numThreads -= 1
		}

		logInfo("All threads completed.")
		endTime := time.Now()
		runDuration := endTime.Sub(startTime)

		logGreen("Completed " + strconv.Itoa(len(uncheckedDomains)) + " domains in " + strconv.FormatFloat(runDuration.Seconds(), 'f', 2, 64) + " seconds")
		logGreen("Total run count: " + strconv.Itoa(totalRuns))
	}

}
