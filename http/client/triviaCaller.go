package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"gihub.com/saiddis/quizgo"
	"github.com/gin-gonic/gin"
)

type TriviaCaller struct {
	client   *http.Client
	Response response
}

func NewTriviaCaller(client *http.Client) *TriviaCaller {
	return &TriviaCaller{
		client:   client,
		Response: response{},
	}
}

type response struct {
	ResponseCode int             `json:"response_code"`
	Results      []quizgo.Trivia `json:"results"`
}

func (t *TriviaCaller) Fetch(c *gin.Context, urls []string) (*[]quizgo.Trivia, error) {

	trivias := []quizgo.Trivia{}
	triviaCaller := t.call(&trivias)
	difficulties := map[int]string{
		0: "medium",
		1: "easy",
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for i, url := range urls {
		//if i < 1 {
		//	time.Sleep(time.Second * 3)
		//}
		go func() {
			defer wg.Done()
			triviaCaller(c, url, difficulties[i])
		}()
	}
	wg.Wait()

	if c.Errors != nil {
		return nil, errors.New("error getting quizzes")
	}
	return &trivias, nil
}

func (t *TriviaCaller) call(trivias *[]quizgo.Trivia) func(c *gin.Context, url, difficulty string) {
	return func(c *gin.Context, url, difficulty string) {
		url += "&" + "difficulty=" + difficulty
		req, err := http.NewRequestWithContext(c, http.MethodGet, url, nil)
		if err != nil {
			c.JSON(400, gin.H{"request error": err.Error()})
			return
		}
		resp, err := t.client.Do(req)
		if err != nil {
			c.JSON(400, gin.H{"response error": err.Error()})
			return
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(400, gin.H{"error reading response body": err.Error()})
			return
		}
		defer resp.Body.Close()

		err = json.Unmarshal(data, &t.Response)
		if err != nil {
			c.JSON(400, gin.H{"error unmarshalling response": err.Error()})
			return
		}

		if t.Response.ResponseCode != 0 {
			c.JSON(400, gin.H{"error": fmt.Sprintf("error fetching data: response code %d", t.Response.ResponseCode)})
			return
		}

		*trivias = append(*trivias, t.Response.Results...)
	}
}
