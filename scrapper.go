package scrapper

import (
	"fmt"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"strings"
	"net/http"
	"io/ioutil"
	"github.com/danigomez/scrapper/util"
)


func mergeMap(current ScrapResult, new ScrapResult, currentUrl string) {
	for k, v := range new {
		newKey := k + "@" + currentUrl
		current[newKey] = v
	}
}

type ScrapResult map[string][]*html.Node

type ScrapDescriptor struct {
	routeSelector string
	tagSelectorMap map[string]string
	next *ScrapDescriptor
}
// TODO Agregar la posibilidad de que se puede definir un RoadMap
// es decir, que se puede indicar que URLs del dominio se deben recorrer a partir de selectores
// y por cada una de esas urls poder definir a su vez nuevos selectores para extraer información
type Scrapper struct {
	domain         string            // Contains the domain where we should start the scrapper
	descriptors []*ScrapDescriptor
}

func NewDescriptor(routeSelector string, tagSelectorMap map[string]string, next *ScrapDescriptor) ScrapDescriptor {
	return ScrapDescriptor{routeSelector, tagSelectorMap, next}
}

func NewScrapper(domain string, descriptor []*ScrapDescriptor) Scrapper {
	return Scrapper{domain, descriptor}
}

func (s Scrapper) DoScrap() (ret ScrapResult) {
	ret = make(ScrapResult)
	response, err := http.Get(s.domain)

	if err != nil {
		fmt.Errorf("error: there was an error while getting resource %s", s.domain)
	}

	body, err := ioutil.ReadAll(response.Body)

	// Iterate over each descriptor
	for _, descriptor := range s.descriptors {
		// For each descriptor, recurse until we reach the end of the linked list
		for descriptor != nil {
			if descriptor.routeSelector == "" {
				mergeMap(ret, s.scrap(string(body), descriptor.tagSelectorMap), s.domain)
			} else {
				aux := s.scrap(string(body), map[string]string{"route": descriptor.routeSelector})
				nodes := aux["route"]
				for _, node := range nodes {
					href := strings.Replace(util.GetValFromKey(node, "href"), "//", "http://", 1)

					if href == "" {
						fmt.Printf("There is no href for selector %s", descriptor.routeSelector)
						continue
					}

					fmt.Printf("Scrapping url %s\n", href)
					response, err = http.Get(href)

					if err != nil {
						fmt.Errorf("error: there was an error while getting resource %s", s.domain)
					}
					body, err = ioutil.ReadAll(response.Body)

					mergeMap(ret, s.scrap(string(body), descriptor.tagSelectorMap), href)
				}

			}

			descriptor = descriptor.next
		}
	}

	return
}

func (s Scrapper) scrap(rawHtml string, tagSelectorMap map[string]string) (ret ScrapResult) {
	// Initialize a new map
	ret = make(ScrapResult)

	// Convert html string to Node representation
	node, err := html.Parse(strings.NewReader(rawHtml))

	if err != nil {
		fmt.Errorf("error: there was and error parsing html %s", err.Error())
	}

	for tag, selector := range tagSelectorMap {

		fmt.Printf("Processing (%s, %s)\n", tag, selector)
		compiled, err := cascadia.Compile(selector)
		if err != nil {
			fmt.Errorf("error: there was an error getting selector %s", err.Error())
		}
		result := compiled.MatchAll(node)

		ret[tag] = result

	}

	return
}
