package economy

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/ninjadotorg/SimEcon/util"
)

const (
	HOST = "http://localhost:8080"
)

func TestSimpleFlow(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	productionId := util.UUID()
	testNewProduction(productionId, `http://localhost:9090/production/firm`, `{"rice":1}`, `{"rice":1}`)
	testProduction(productionId)

	agentId1 := util.UUID()
	testNewAgent(agentId1, productionId)
	testAgent(agentId1)

	assetId := "rice"
	testNewMarket(assetId)
	testMarket(assetId)

	testBuyLimit(assetId, "10", "1", agentId1)
	testMarket(assetId)

	testBuyLimit(assetId, "20", "1", agentId1)
	testMarket(assetId)

	testBuyLimit(assetId, "30", "0.9", agentId1)
	testMarket(assetId)

	testBuyLimit(assetId, "40", "1.1", agentId1)
	testMarket(assetId)

	agentId2 := util.UUID()
	testNewAgent(agentId2, productionId)
	testAgent(agentId2)

	testSell(assetId, "10", agentId2)
	testMarket(assetId)
	testAgent(agentId1)
	testAgent(agentId2)

	testSellLimit(assetId, "10", "1.3", agentId2)
	testMarket(assetId)

	testSellLimit(assetId, "20", "1.3", agentId2)
	testMarket(assetId)

	testSellLimit(assetId, "30", "1.4", agentId2)
	testMarket(assetId)

	testSellLimit(assetId, "40", "1.2", agentId2)
	testMarket(assetId)

	testAgentProduce(agentId1, `{"rice":1}`)

}

func testAgentProduce(agentId, input string) {
	testGet(HOST + "/agent/" + agentId + "/produce" +
		"?input=" + input)
}

func testMarket(assetId string) {
	testGet(HOST + "/market/" + assetId)
}

func testAgent(agentId string) {
	testGet(HOST + "/agent/" + agentId)
}

func testProduction(productionId string) {
	testGet(HOST + "/production/" + productionId)
}

func testSellLimit(assetId, size, price, agentId string) {
	testGet(HOST + "/market/" + assetId + "/sellLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId)
}

func testBuyLimit(assetId, size, price, agentId string) {
	testGet(HOST + "/market/" + assetId + "/buyLimit" +
		"?size=" + size +
		"&price=" + price +
		"&agentId=" + agentId)
}

func testSell(assetId, amount, agentId string) {
	testGet(HOST + "/market/" + assetId + "/sell" +
		"?amount=" + amount +
		"&agentId=" + agentId)
}

func testNewAgent(agentId, productionId string) {
	testGet(HOST + "/agent/" + agentId + "/new" +
		"?productionId=" + productionId)
}

func testNewMarket(assetId string) {
	testGet(HOST + "/market/" + assetId + "/new")
}

func testNewProduction(productionId, function, input, output string) {
	testGet(HOST + "/production/" + productionId + "/new" +
		"?function=" + function +
		"&input=" + input +
		"&output=" + output)
}

func testGet(url string) {
	log.Println("-----")
	log.Println("REQ", url)
	if res, e := http.Get(url); e == nil {
		if data, e := ioutil.ReadAll(res.Body); e == nil {
			log.Println("RES", string(data))
		} else {
			log.Println("ERR", e)
		}
	} else {
		log.Println("ERR", e)
	}
}
