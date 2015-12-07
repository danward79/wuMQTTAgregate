package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/danward79/mqttservices"
	"github.com/danward79/sensorCache"
	"github.com/danward79/wupws"
	proto "github.com/huin/mqtt"
)

var mqttServer, password, stationID, software, configPath *string
var calculateDewpoint *bool
var config map[string]string
var addressParameter map[string]string
var sensorExpire, stationReportPeriod, checkCache *time.Duration
var wg sync.WaitGroup
var done chan struct{}

func init() {
	//Command line variables
	mqttServer = flag.String("s", ":1883", "IP and Port of the MQTT Broker. e.g. 127.0.0.1:1883. Default: :1883")
	stationID = flag.String("u", "", "WU PWS station id")
	password = flag.String("p", "", "WU PWS station password")
	software = flag.String("f", "gowupws", "Name of the software updating the PWS. Default: gowupws")
	calculateDewpoint = flag.Bool("d", false, "Provide calculated dewpoint, if not provided as a parameter. Default: False")
	configPath = flag.String("c", "", "Provide path to config file")
	sensorExpire = flag.Duration("l", 5*time.Minute, "Sensor/Device life, minutes")
	checkCache = flag.Duration("e", 1*time.Minute, "Check sensor/device every, minutes")
	stationReportPeriod = flag.Duration("r", 2*time.Minute, "Station report period, minutes")
	flag.Parse()

	if *stationID == "" {
		log.Fatal("A Weather Underground station ID has to be provided")
	}

	if *password == "" {
		log.Fatal("A Weather Underground password has to be provided")
	}

	if *configPath == "" {
		log.Fatal("A config file has to be provided")
	}

	config = readConfigFile(*configPath)
	addressParameter = mapAddressToParameter(&config)

	done = make(chan struct{})
}

func main() {
	//PWS - Personal Weather Station
	pws := wupws.New(*stationID, *password, *software, *calculateDewpoint)
	mqttClient := mqttservices.NewClient(*mqttServer)
	cache := sensorCache.New(*sensorExpire)

	//go routine to check if the data stored in the cache is expired
	go cache.MonitorExpiry(*checkCache)

	//go routine to capture readings into cache
	wg.Add(1)
	go cacheReadings(cache, mqttClient)

	//go routine to do the updates
	wg.Add(1)
	go pushUpdates(cache, *stationReportPeriod, pws)

	wg.Wait()

	log.Println("wuMQTTAgregate: Exiting")
}

//subscribeSensors uses config data to setup subscriptions to sensors that feed the WU API
func subscribeSensors(c *map[string]string, m *mqttservices.MqttClient) chan *proto.Publish {
	topicList := make([]proto.TopicQos, len(config))
	i := 0
	for _, v := range config {
		topicList[i].Topic = v
		topicList[i].Qos = proto.QosAtMostOnce
		i++
	}

	return m.Subscribe(topicList)
}

//mapAddressToParameter
func mapAddressToParameter(c *map[string]string) map[string]string {
	m := make(map[string]string)
	for k, v := range *c {
		m[v] = k
	}
	return m
}

//cacheReadings catch readings from subscribed channels and cache
func cacheReadings(c *sensorCache.Cache, s *mqttservices.MqttClient) {

	defer wg.Done()

	chIn := subscribeSensors(&config, s)

	for m := range chIn {
		c.Insert(addressParameter[m.TopicName], fmt.Sprintf("%s", m.Payload))
	}

	log.Println("wuMQTTAgregate: MQTT broker connection closed")
	c.StopMonitoring()
	close(done)
}

//getCacheReadings grabs the readings and returns a map[string]string
func getCacheReadings(c *sensorCache.Cache) map[string]string {

	m := make(map[string]string)
	for k, v := range c.Values() {
		if _, OK := v.(string); OK {
			m[k] = v.(string)
		}
	}
	return m
}

//pushUpdates to wu
func pushUpdates(c *sensorCache.Cache, td time.Duration, stn *wupws.Station) {
	defer wg.Done()

	//continuous loop to report PWS data
	t := time.NewTicker(td)

	for {
		select {
		case <-t.C:
			stn.UpdateWeather(getCacheReadings(c))

			err := stn.PushUpdate("")
			if err != nil {
				log.Println(err)
			}

		case <-done:
			log.Println("wuMQTTAgregate: Pushing stopped")
			return
		}
	}

}
