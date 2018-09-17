package economy

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/ninjadotorg/SimEcon/util"
)

func TestSimpleFlow(t *testing.T) {

	// add production
	productionId := util.UUID()
	log.Println("productionId", productionId)

	function := `http://localhost:9090/production/firm`
	input := `{"rice":1}`
	output := `{"meal":1}`

	url := "http://localhost:8080/production/" + productionId + "/new" +
		"?function=" + function +
		"&input=" + input +
		"&output=" + output
	log.Println("REQUEST", url)
	res, _ := http.Get(url)
	data, _ := ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view production
	url = "http://localhost:8080/production/" + productionId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// add agent
	agentId1 := util.UUID()
	log.Println("agentId", agentId1)

	url = "http://localhost:8080/agent/" + agentId1 + "/new" +
		"?productionId=" + productionId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view agent
	url = "http://localhost:8080/agent/" + agentId1
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// add market
	assetId := "rice"
	log.Println("assetId", assetId)

	url = "http://localhost:8080/market/" + assetId + "/new"
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// first buy limit order
	assetId = "rice"
	size := "10"
	price := "1"
	url = "http://localhost:8080/market/" + assetId + "/buyLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId1

	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view market
	url = "http://localhost:8080/market/" + assetId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// second buy limit order, same price
	assetId = "rice"
	size = "20"
	price = "1"
	url = "http://localhost:8080/market/" + assetId + "/buyLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId1

	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view market
	url = "http://localhost:8080/market/" + assetId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// third buy limit order, lower price
	assetId = "rice"
	size = "30"
	price = "0.9"
	url = "http://localhost:8080/market/" + assetId + "/buyLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId1

	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view market
	url = "http://localhost:8080/market/" + assetId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// fourth buy limit order, higher price
	assetId = "rice"
	size = "40"
	price = "1.1"
	url = "http://localhost:8080/market/" + assetId + "/buyLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId1

	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view market
	url = "http://localhost:8080/market/" + assetId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// add agent
	agentId2 := util.UUID()
	log.Println("agentId2", agentId2)

	url = "http://localhost:8080/agent/" + agentId2 + "/new" +
		"?productionId=" + productionId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// first sell market order
	assetId = "rice"
	amount := "10"
	url = "http://localhost:8080/market/" + assetId + "/sell" +
		"?amount=" + amount +
		"&agentId=" + agentId2

	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view market
	url = "http://localhost:8080/market/" + assetId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// first sell limit order
	assetId = "rice"
	size = "10"
	price = "1.3"
	url = "http://localhost:8080/market/" + assetId + "/sellLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId2

	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view market
	url = "http://localhost:8080/market/" + assetId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// second sell limit order, same price
	assetId = "rice"
	size = "20"
	price = "1.3"
	url = "http://localhost:8080/market/" + assetId + "/sellLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId2

	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view market
	url = "http://localhost:8080/market/" + assetId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// third sell limit order, higher price
	assetId = "rice"
	size = "30"
	price = "1.4"
	url = "http://localhost:8080/market/" + assetId + "/sellLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId2

	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view market
	url = "http://localhost:8080/market/" + assetId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// fourth sell limit order, lower price
	assetId = "rice"
	size = "40"
	price = "1.2"
	url = "http://localhost:8080/market/" + assetId + "/sellLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId2

	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))

	// view market
	url = "http://localhost:8080/market/" + assetId
	log.Println("REQUEST", url)
	res, _ = http.Get(url)
	data, _ = ioutil.ReadAll(res.Body)
	log.Println("RESPONSE", string(data))
}
