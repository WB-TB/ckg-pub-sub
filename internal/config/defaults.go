package config

func setDefaults() map[string]any {
	return map[string]any{
		// App
		"app.env":      "development",
		"app.loglevel": "info",

		// Google Cloud
		"google.project":     "ckg-tb-staging",
		"google.credentials": "./credentials.json",
		"google.debug":       false,

		// PubSub
		"pubsub.topic":           "projects/ckg-tb-staging/topics/CKG-SITB",
		"pubsub.subscription":    "projects/ckg-tb-staging/subscriptions/CKG-SITB-sub",
		"pubsub.messageordering": false,

		// Consumer
		"consumer.maxmessages":             10,
		"consumer.sleeptime":               "5s",
		"consumer.acktimeout":              "60s",
		"consumer.retrycount":              3,
		"consumer.retrydelay":              "1s",
		"consumer.flowcontrol.enabled":     true,
		"consumer.flowcontrol.maxmessages": 1000,
		"consumer.flowcontrol.maxbytes":    1000000, // 1M

		// Producer
		"producer.enableordering":        false,
		"producer.batchsize":             100,
		"producer.compression.enabled":   false,
		"producer.compression.algorithm": "gzip",

		// API
		"api.baseurl":   "https://api-dev.dto.kemkes.go.id/fhir-sirs",
		"api.timeout":   "60s",
		"api.apiheader": "X-API-Key:",
		"api.batchsize": 100,

		// Database
		"db.driver": "mongodb",
		"db.host":   "localhost",
		"db.port":   27017,
		// "db.username":   "xtb",
		// "db.password":   "xtb",
		"db.database":   "ckgtb",
		"db.attributes": "",

		// CKG
		"ckg.usecache":           false,
		"ckg.tablemasterwilayah": "master_wilayah",
		"ckg.tablemasterfaskes":  "master_faskes",
		"ckg.tableskrining":      "skrining_tb",
		"ckg.tablestatus":        "pasien_tb",
		"ckg.tableincoming":      "ckg_pubsub_incoming",
		"ckg.tableoutgoing":      "ckg_pubsub_outgoing",
		"ckg.markerfield":        "transactionSource",
		"ckg.markerconsume":      "STATUS-PASIEN-TB",
		"ckg.markerproduce":      "SKRINING-CKG-TB",
	}
}
