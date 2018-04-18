package util

import "golang.org/x/net/html"


func GetValFromKey(node *html.Node, key string) string {
	for _, v := range node.Attr {
		if v.Key == key {
			return v.Val
		}
	}

	return ""
}
