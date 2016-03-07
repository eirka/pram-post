package middleware

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/eirka/eirka-libs/config"
)

var (
	ErrBlacklisted = errors.New("IP is on spam blacklist")
)

// check ip with stop forum spam
func StopSpam() gin.HandlerFunc {
	return func(c *gin.Context) {

		// check ip against stop forum spam
		err := CheckStopForumSpam(c.ClientIP())
		if err == ErrBlacklisted {
			c.JSON(http.StatusBadRequest, gin.H{"error_message": "IP is on spam blacklist"})
			c.Error(err).SetMeta("StopSpam.CheckStopForumSpam")
			c.Abort()
			return
		}

		c.Next()

	}
}

// Stop Forum Spam return format
type StopForumSpam struct {
	Ip struct {
		Appears    float64 `json:"appears"`
		Confidence float64 `json:"confidence"`
		Frequency  float64 `json:"frequency"`
		Lastseen   string  `json:"lastseen"`
	} `json:"ip"`
	Success float64 `json:"success"`
}

// Check Stop Forum Spam blacklist for IP
func CheckStopForumSpam(ip string) (err error) {

	if len(ip) == 0 {
		return errors.New("no ip provided")
	}

	queryValues := url.Values{}

	queryValues.Set("ip", ip)
	queryValues.Set("f", "json")

	// construct the api request
	sfs_endpoint := &url.URL{
		Scheme:   "http",
		Host:     "api.stopforumspam.org",
		Path:     "api",
		RawQuery: queryValues.Encode(),
	}

	// our http request
	req, err := http.NewRequest(http.MethodGet, sfs_endpoint.String(), nil)
	if err != nil {
		return errors.New("Error creating SFS request")
	}

	// set ua header
	req.Header.Set("User-Agent", "Eirka/1.2")

	// a client with a timeout
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	// do the request
	// TODO: add errors here to a system log
	resp, err := netClient.Do(req)
	if err != nil {
		return errors.New("Error reaching SFS")
	}
	defer resp.Body.Close()

	// read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error parsing SFS response")
	}

	sfs_data := StopForumSpam{}

	// unmarshal into struct
	err = json.Unmarshal(body, &sfs_data)
	if err != nil {
		return errors.New("Error parsing SFS data")
	}

	// check if the spammer confidence level is over our setting
	if sfs_data.Ip.Confidence > config.Settings.StopForumSpam.Confidence {
		return ErrBlacklisted
	}

	return

}
