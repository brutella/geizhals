package geizhals

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ericchiang/css"
	"golang.org/x/net/html"
)

type Domain string

const (
	DomainAt Domain = "geizhals.at"
	DomainEu Domain = "geizhals.eu"
	DomainDe Domain = "geizhals.de"
	DomainUk Domain = "skinflint.co.uk"
	DomainPl Domain = "cenowarka.pl"
)

type Product struct {
	Id       string
	Name     string
	MinPrice float64
	MaxPrice float64
}

func GetProduct(id string) (*Product, error) {
	return DefaultClient.Get(id)
}

type Client struct {
	Timeout time.Duration
	Domain  Domain
}

var DefaultClient = &Client{
	Timeout: 2 * time.Second,
	Domain:  DomainEu,
}

func (cl *Client) Get(id string) (*Product, error) {
	url := fmt.Sprintf("https://%s/%s", cl.Domain, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return cl.parseProduct(id, resp.Body)
}

func (cl *Client) parseProduct(id string, reader io.Reader) (*Product, error) {
	node, err := html.Parse(reader)
	if err != nil {
		panic(err)
	}

	min, max, err := cl.parsePrice(node)
	if err != nil {
		return nil, err
	}

	name, err := cl.parseName(node)
	if err != nil {
		return nil, err
	}

	pr := &Product{
		Id:       id,
		Name:     name,
		MinPrice: min,
		MaxPrice: max,
	}
	return pr, nil
}

func (cl *Client) parseName(node *html.Node) (string, error) {
	sel, err := css.Compile(".variant__header__headline")
	if err != nil {
		panic(err)
	}
	for _, ele := range sel.Select(node) {
		if child := ele.FirstChild; child != nil {
			return strings.Trim(child.Data, " \n\r"), nil
		}
	}

	return "", fmt.Errorf("product name not found")
}

func (cl *Client) parsePrice(node *html.Node) (min float64, max float64, err error) {
	sel, err := css.Compile(".variant__header__pricehistory__pricerange .gh_price")
	if err != nil {
		panic(err)
	}
	for _, ele := range sel.Select(node) {
		if child := ele.FirstChild; child != nil {
			var val float64
			if val, err = cl.parseFloat(child.Data); err != nil {
				return
			} else if min == 0 {
				min = val
			} else {
				max = val
			}
		}
	}

	return
}

func (cl *Client) parseFloat(str string) (float64, error) {
	switch cl.Domain {
	case DomainAt, DomainEu, DomainDe:
		return cl.parseEuro(str)
	default:
		return 0, fmt.Errorf("Unable to parse amount %s", str)
	}
}

func (cl *Client) parseEuro(str string) (float64, error) {
	str = strings.Replace(str, "€", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, ".", "", -1)
	str = strings.Replace(str, ",", ".", -1)
	return strconv.ParseFloat(str, 64)
}
