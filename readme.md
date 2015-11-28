# WuMQTTAgregate

With many climate related parameters reporting as topics to a MQTT server, this service subscribes to *MQTT Topics*, which are related to weather and periodically pushed the parameters to [Weather Underground](www.weatherunderground.com).

Weather Underground provide a simple web HTTP GET API. See [here](http://wiki.wunderground.com/index.php/PWS_-_Upload_Protocol)

In order to provide this functionality, this app uses two libraries. wupws, which provides an interface to Weather Underground and sensorcache, which provides and manages a cache of sensor readings with automatic expiry.
