package ddg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/matrix-org/go-neb/types"
	"github.com/matrix-org/gomatrix"
	"github.com/puerkitobio/goquery"
	log "github.com/sirupsen/logrus"
)

const ServiceType = "ddg"

const searchEndpoint = "https://duckduckgo.com/html"
const defaultUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36"

var httpClient = &http.Client{}

type Service struct {
	types.DefaultService
}

func (s *Service) Commands(cli *gomatrix.Client) []types.Command {
	return []types.Command{
		types.Command{
			Path:    []string{"g"},
			Command: search,
		},
		types.Command{
			Path:    []string{"п"},
			Command: search,
		},
	}
}

func search(roomID, userID string, args []string) (interface{}, error) {
	if len(args) == 0 {
		return usage(), nil
	}
	search := strings.Join(args, " ")

	req, _ := http.NewRequest("GET", searchEndpoint, nil)
	req.Header.Set("User-Agent", defaultUA)
	q := req.URL.Query()
	q.Add("q", strings.Join(args, " "))
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 200 {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, responseToString(resp))
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return notFound(search), nil
	}
	doc.Find(".result--ad").Remove()
	link, ok := doc.Find(".result__a").First().Attr("href")
	if !ok {
		return notFound(search), nil
	}
	return found(search, link), nil
}

func usage() *gomatrix.TextMessage {
	return &gomatrix.TextMessage{"m.notice", "Использование: !g search_text"}
}

func notFound(search string) *gomatrix.TextMessage {
	return &gomatrix.TextMessage{"m.notice", fmt.Sprintf("Я ничего не нашел про %s.", search)}
}

func found(search, link string) *gomatrix.TextMessage {
	return &gomatrix.TextMessage{"m.notice", fmt.Sprintf("Я нашел %s", link)}
}

func responseToString(resp *http.Response) string {
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Failed to decode response body"
	}
	return string(bs)
}

func init() {
	types.RegisterService(func(serviceID, serviceUserID, arg3 string) types.Service {
		return &Service{types.NewDefaultService(serviceID, serviceUserID, ServiceType)}
	})
	log.Info("ddg service loaded")
}
