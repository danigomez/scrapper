package scrapper

import (
	"testing"
	"fmt"
	"strings"
)

func TestScrap(t *testing.T) {

	hasKey := false
	data := `
		<body>
			<div class="link">
				<a href="www.google.com"> Test </a>
			</div>
			<div class="other_link">
				<a href="www.facebook.com"> Test 2 </a>
			</div>
		</body>
	`
	scrapper := NewScrapper("", nil)

	ret := scrapper.scrap(data, map[string]string{"test": "div.link a[href]"})

	if len(ret) > 1 {
		t.Errorf("The matched elements are more than 1 => %v", len(ret))
	}


	for k := range ret {
		if k == "test" {
			hasKey = true
			break
		}
	}

	if !hasKey {
		t.Errorf("The key 'test' doesn't exists!")
	}

	first := ret["test"][0]

	if len(first.Attr) > 1 {
		t.Errorf("The data matched should have 1 attribute => %s", first.Attr)
	}
	for _, v := range first.Attr {
		if v.Key == "href" {
			if v.Val != "www.google.com" {
				t.Errorf("href value of matched element should be www.google.com => %s", v.Val)
			}
			fmt.Printf("Link => %s\n", v.Val)

			break
		}
	}

}

func TestDoScrapWithoutRouteSelector(t *testing.T)  {
	tags := map[string]string{
		"googleImg": "div#lga img#hplogo[src]",
	}
	descriptor := NewDescriptor("", tags, nil)
	scrapper := NewScrapper(
		"http://www.google.com",
		[]*ScrapDescriptor{
			&descriptor,
		})
	ret := scrapper.DoScrap()

	img := ret["googleImg@http://www.google.com"][0]

	if img == nil {
		t.Errorf("There was no img tag for Google")
	}

	for _, v := range img.Attr {
		if v.Key == "src" {
			if v.Val == "" {
				t.Errorf("There is no src for image")
			}
			fmt.Printf("Google image Link => %s%s\n", scrapper.domain, v.Val)
			break
		}
	}


}


func TestDoScrapWithRouteSelector(t *testing.T)  {
	tags := map[string]string{
		"productTitle": "div.gb-list-cluster h3.gb-list-cluster-title",
	}
	descriptor := NewDescriptor("div.gb-category-submenu-title > a[href]", tags, nil)
	scrapper := NewScrapper(
		"http://www.garbarino.com",
		[]*ScrapDescriptor{
			&descriptor,
		})
	ret := scrapper.DoScrap()

	if ret == nil {
		t.Errorf("There was no result")
	}

	for key, title := range ret {
		fmt.Printf("\n\nShowing key => %s\n", key)
		for _, node := range title {
			child := node.FirstChild
			fmt.Println(strings.TrimSpace(child.Data))
		}
	}




}