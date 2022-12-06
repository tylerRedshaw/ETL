package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antihax/optional"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sqconnect "github.com/square/square-connect-go-sdk/swagger"
)

var (
	cfg = &sqconnect.Configuration{
		BasePath:      "https://connect.squareupsandbox.com",
		DefaultHeader: make(map[string]string),
		UserAgent:     "Swagger-Codegen/1.0.0/go",
	}
	fixedPricingType                 = sqconnect.FIXED_PRICING_CatalogPricingType
	usd_currency                     = sqconnect.USD_Currency
	ITEM_VARIATION_CatalogObjectType = sqconnect.ITEM_VARIATION_CatalogObjectType
)

func init() {
	fmt.Println("This will get called on main initialization")
	cfg.AddDefaultHeader("Authorization", "Bearer EAAAEINRS1-ATVKx_ZBs2oVgffIzRtDcDcJ7LReuJTdCs4Qo1ECr7yCmwQPgPRJr")
}

type BatchUpdate struct {
	CatalogObjectBatch sqconnect.CatalogObjectBatch `json:"catalogObjectBatch"`
	ApiKey             string                       `json:"apiKey"`
}

// using go SDK
func ListCatalog(c *gin.Context) {
	client := sqconnect.NewAPIClient(cfg)
	ctx := context.TODO()
	listCatalogResponse, _, err := client.CatalogApi.ListCatalog(ctx, &sqconnect.CatalogApiListCatalogOpts{})

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println("BODY STRING")
	fmt.Println(bodyString)

	// return data
	fmt.Println(listCatalogResponse, err)
	c.IndentedJSON(http.StatusOK, listCatalogResponse)
}

func UpdateCatalogObject(c *gin.Context) {
	client := sqconnect.NewAPIClient(cfg)
	fmt.Println("REQUEST BODY")
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println("BODY STRING")
	fmt.Println(bodyString)

	ctx := context.TODO()
	upsertCatalogObjectRequest := sqconnect.UpsertCatalogObjectRequest{
		IdempotencyKey: uuid.New().String(),
		Object: &sqconnect.CatalogObject{
			Id:    "QTSNJYCWGNBSGAFZN2FSQQ5J",
			Type_: &ITEM_VARIATION_CatalogObjectType,
			ItemVariationData: &sqconnect.CatalogItemVariation{
				Sku:         "123456789",
				PricingType: &fixedPricingType,
				PriceMoney: &sqconnect.Money{
					Amount:   800,
					Currency: &usd_currency,
				},
				ItemId: "HTDOPMSWWQ22YMBNDHW6ICCH",
				Name:   "Blah2",
			},
			Version: 1669990656318,
		},
	}
	upsertCatalogObjectResponse, res, err := client.CatalogApi.UpsertCatalogObject(ctx, upsertCatalogObjectRequest)

	if err == nil {
		fmt.Println("Successful Catalog Upsert")
		fmt.Println(upsertCatalogObjectResponse)
		c.IndentedJSON(http.StatusOK, upsertCatalogObjectResponse)
	} else {
		fmt.Println(res)
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, nil)
	}
}

func UpdateBatchCatalogObject(c *gin.Context) {
	client := sqconnect.NewAPIClient(cfg)
	fmt.Println("REQUEST BODY")
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println("BODY STRING")
	fmt.Println(bodyString)

	var catalogObjectsBatch BatchUpdate
	json.Unmarshal([]byte(bodyBytes), &catalogObjectsBatch)
	fmt.Println("TESTING API KEY:      ", catalogObjectsBatch.ApiKey)

	batchUpsertCatalogObjectsRequest := sqconnect.BatchUpsertCatalogObjectsRequest{
		IdempotencyKey: uuid.New().String(),
		Batches:        []sqconnect.CatalogObjectBatch{catalogObjectsBatch.CatalogObjectBatch},
	}
	ctx := context.TODO()

	upsertCatalogObjectResponse, res, err := client.CatalogApi.BatchUpsertCatalogObjects(ctx, batchUpsertCatalogObjectsRequest)

	if err == nil {
		fmt.Println("Successful Catalog Upsert")
		fmt.Println(upsertCatalogObjectResponse)
		c.IndentedJSON(http.StatusOK, upsertCatalogObjectResponse)
	} else {
		fmt.Println(res)
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, nil)
	}
}

// without using the go sdk
func OrderSimple(c *gin.Context) {
	url := "https://connect.squareup.com/v2/orders"
	method := "POST"

	payload := strings.NewReader(`{
    "idempotency_key": "10d03b4f-7641-4f2a-8812-0196caa6a600",
    "order": {
      "location_id": "LADNP759ZA8PC",
      "state": "OPEN",
      "service_charges": [
        {
          "amount_money": {
            "amount": 99,
            "currency": "USD"
          },
          "name": "Test Order",
          "calculation_phase": "TOTAL_PHASE"
        }
      ]
    }
  }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Square-Version", "2022-11-16")
	req.Header.Add("Authorization", "Bearer EAAAEXFV8pPcWCE2dexBcdgsLv-qLZ3L93eUg4J1NrBX12GTkY9lPTfZ1uLPybF5")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// using go sdk for json marshalling
	var createOrderReponse sqconnect.CreateOrderResponse
	json.Unmarshal([]byte(body), &createOrderReponse)
	fmt.Printf("Id: %s, LocationId: %s", createOrderReponse.Order.Id, createOrderReponse.Order.LocationId)
	fmt.Println(string(body))
	c.IndentedJSON(http.StatusOK, createOrderReponse)
}

func ListCustomers(c *gin.Context) {
	// TODO(tredshaw): add support for API Key from app
	client := sqconnect.NewAPIClient(cfg)
	ctx := context.TODO()

	var sortOrder interface{} = "DESC"
	var sortField interface{} = "CREATED_AT"

	listCustomersResponse, _, err := client.CustomersApi.ListCustomers(ctx, &sqconnect.CustomersApiListCustomersOpts{
		SortOrder: optional.NewInterface(sortOrder),
		SortField: optional.NewInterface(sortField),
	})

	// return data
	fmt.Println(listCustomersResponse, err)
	c.IndentedJSON(http.StatusOK, listCustomersResponse)
}
