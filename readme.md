# WuMQTTAgregate

With many climate related parameters reporting as topics to a MQTT server, this service subscribes to *MQTT Topics*, which are related to weather and periodically pushed the parameters to [Weather Underground](www.weatherunderground.com).

Weather Underground provide a simple web HTTP GET API. See [here](http://wiki.wunderground.com/index.php/PWS_-_Upload_Protocol)

In order to provide this functionality, this app uses two libraries. [wupws](github.com/danward79/wupws), which provides an interface to Weather Underground and [sensorcache](github.com/danward79/sensorcache), which provides and manages a cache of sensor readings with automatic expiry.

## Install

```
go get -u github.com/danward79/WuMQTTAgregate
go install
```

## Run

In order to use this code you need to have a personal weather station at WeatherUnderground set up.

```
wuMQTTAgregate -s 192.168.0.7:1883 -u STNID -p PASSWORD -c ./config.cfg
```
