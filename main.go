package main

import (
	// "fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"log"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var jobs []string
	sites := map[string]string{"https://www.bund.net/ueber-uns/jobs": ".m-content-faq--header", "http://wwf.panda.org/jobs_wwf": ".media-heading a", "https://www.boell.de/de/jobs-der-heinrich-boell-stiftung": ".teaser__title a"}
	for k, v := range sites {
		jobs = append(jobs, postScrape(k, v)...)
	}
	// fmt.Print(jobs)
	// fmt.Print(len(jobs))
	mailer(jobs)

}
func postScrape(site, job string) []string {
	doc, err := goquery.NewDocument(site)
	if err != nil {
		log.Fatal(err)
	}

	var hits []string

	doc.Find(job).Each(func(index int, item *goquery.Selection) {
		jobs := []string{"Presse", "Communications", "Forest", "Feminismus"}
		title := strings.TrimSpace(item.Text())
		foundJob := contains(jobs, title)
		if foundJob {
			// linkTag := item.Find("a")
			// link, _ := linkTag.Attr("href")
			// fmt.Printf("index %d\n", index)
			// fmt.Printf("title %s\n", title)
			// fmt.Printf("link %s\n\n", site+link)
			// mailBody.WriteString(title + "\n")
			hits = append(hits, title)
		}
	})
	return hits
}
func contains(arr []string, str string) bool {
	for _, a := range arr {
		i := strings.Index(str, a)
		if i >= 0 {
			return true
		} else {
			// fmt.Printf("Nothing found in %s\n", str)
		}
	}
	return false
}
