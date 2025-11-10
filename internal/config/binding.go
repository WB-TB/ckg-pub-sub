package config

func setEnvBindings() map[string]string {
	return map[string]string{
		// App
		"app.env":      "APP_ENV",
		"app.loglevel": "APP_LOGLEVEL",

		// Google Cloud
		"google.project":     "GOOGLE_PROJECT",
		"google.credentials": "GOOGLE_CREDENTIALS",
		"google.debug":       "GOOGLE_DEBUG",

		// PubSub
		"pubsub.topic":           "PUBSUB_TOPIC",
		"pubsub.subscription":    "PUBSUB_SUBSCRIPTION",
		"pubsub.messageordering": "PUBSUB_MESSAGEORDERING",

		// Consumer
		"consumer.maxmessages":             "CONSUMER_MAXMESSAGES",
		"consumer.sleeptime":               "CONSUMER_SLEEPTIME",
		"consumer.acktimeout":              "CONSUMER_ACKTIMEOUT",
		"consumer.retrycount":              "CONSUMER_RETRYCOUNT",
		"consumer.retrydelay":              "CONSUMER_RETRYDELAY",
		"consumer.flowcontrol.enabled":     "CONSUMER_FLOWCONTROL_ENABLED",
		"consumer.flowcontrol.maxmessages": "CONSUMER_FLOWCONTROL_MAXMESSAGES",
		"consumer.flowcontrol.maxbytes":    "CONSUMER_FLOWCONTROL_MAXBYTES",

		// Producer
		"producer.enableordering":        "PRODUCER_ENABLEORDERING",
		"producer.batchsize":             "PRODUCER_BATCHSIZE",
		"producer.compression.enabled":   "PRODUCER_COMPRESSION_ENABLED",
		"producer.compression.algorithm": "PRODUCER_COMPRESSION_ALGORITHM",

		// API
		"api.baseurl":   "API_BASEURL",
		"api.timeout":   "API_TIMEOUT",
		"api.apikey":    "API_APIKEY",
		"api.apiheader": "API_APIHEADER",
		"api.batchsize": "API_BATCHSIZE",

		// Database
		"db.driver":     "DB_DRIVER",
		"db.host":       "DB_HOST",
		"db.port":       "DB_PORT",
		"db.username":   "DB_USERNAME",
		"db.password":   "DB_PASSWORD",
		"db.database":   "DB_DATABASE",
		"db.attributes": "DB_ATTRIBUTES",

		// CKG
		"ckg.tableskrining": "CKG_TABLESKRINING",
		"ckg.tablestatus":   "CKG_TABLESTATUS",
		"ckg.tableincoming": "CKG_TABLEINCOMING",
		"ckg.tableoutgoing": "CKG_TABLEOUTGOING",
		"ckg.markerfield":   "CKG_MARKERFIELD",
		"ckg.markerconsume": "CKG_MARKERCONSUME",
		"ckg.markerproduce": "CKG_MARKERPRODUCE",
	}
}
