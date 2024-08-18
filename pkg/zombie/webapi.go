package zombie

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	apiHeaderSecret    = "X-SDTD-API-SECRET"
	apiHeaderTokenName = "X-SDTD-API-TOKENNAME"
)

type API struct {
	client         *http.Client
	Addr           string
	ApiTokenName   string
	ApiTokenSecret string
}

func NewAPI() *API {
	return &API{
		Addr:           cmp.Or(os.Getenv("SDTD_API_ADDR"), "https://7days.burpcraft.com"),
		ApiTokenName:   os.Getenv("STDT_API_TOKEN_NAME"),
		ApiTokenSecret: os.Getenv("STDT_API_TOKEN_SECRET"),
		client:         cleanhttp.DefaultClient(),
	}
}

type ServerResponse[T any] struct {
	Data T              `json:"data"`
	Meta map[string]any `json:"meta"`
}

type ServerStats struct {
	GameTime struct {
		Days    int `json:"days"`
		Hours   int `json:"hours"`
		Minutes int `json:"minutes"`
	} `json:"gameTime"`
	Players   int       `json:"players"`
	Hostiles  int       `json:"hostiles"`
	Animals   int       `json:"animals"`
	Timestamp time.Time `json:"-"`
}

func (s ServerStats) String() string {
	return fmt.Sprintf("Day %d, %02d:%02d - Players: %d, Hostiles: %d, Animals: %d",
		s.GameTime.Days, s.GameTime.Hours, s.GameTime.Minutes, s.Players, s.Hostiles, s.Animals)
}

func (s ServerStats) Duration() time.Duration {
	return time.Duration(s.GameTime.Days*24*60*60+s.GameTime.Hours*60*60+s.GameTime.Minutes*60) * time.Second
}

func (a *API) ServerStats() (ServerStats, error) {
	var response ServerResponse[ServerStats]
	err := a.Get("/api/serverstats", &response)
	response.Data.Timestamp = time.Now()
	return response.Data, err
}

func (a *API) Get(path string, response any) error {
	req, err := http.NewRequest("GET", a.Addr+path, nil)
	if err != nil {
		return err
	}

	if a.ApiTokenName != "" && a.ApiTokenSecret != "" {
		req.Header.Set(apiHeaderTokenName, a.ApiTokenName)
		req.Header.Set(apiHeaderSecret, a.ApiTokenSecret)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}

	return nil
}
