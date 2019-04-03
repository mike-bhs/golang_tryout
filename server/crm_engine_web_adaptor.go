package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mike-bhs/golang_tryout/app/models"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (serv *Server) GetAllIbans(c *gin.Context) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", "http://localhost:38081/crm/ibans/external/find", nil)
	addQueryParameters(req)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"reason": "failed to perform request",
		})

		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"reason": "Failed to parse body",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"response": parseBody(bodyBytes),
	})
}

func addQueryParameters(req *http.Request) {
	q := req.URL.Query()

	q.Add("page", "1")
	q.Add("per_page", "3")

	q.Add("principal_identifier[contact_id]", "58e78791-e0e5-012c-2dee-001e52f3c730")
	q.Add("principal_identifier[account_id]", "2090939e-b2f7-3f2b-1363-4d235b3f58af")
	q.Add("principal_identifier[effective_contact_id]", "58e78791-e0e5-012c-2dee-001e52f3c730")
	q.Add("principal_identifier[effective_account_id]", "2090939e-b2f7-3f2b-1363-4d235b3f58af")

	req.URL.RawQuery = q.Encode()
}

func parseBody(body []byte) models.Ibans {
	var ibans models.Ibans
	var iban models.Iban

	ibansJson := gjson.Get(string(body), "data.ibans").Array()

	for _, res := range ibansJson {
		err := json.Unmarshal([]byte(res.String()), &iban)

		if err != nil {
			log.Println("Failed to prase json for iban", res.Get("uuid"))
			continue
		}

		ibans = append(ibans, &iban)
	}

	return ibans
}

