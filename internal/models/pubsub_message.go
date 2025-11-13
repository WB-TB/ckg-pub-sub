package models

type IncomingMessageStatusTB struct {
	ID          string  `json:"id" bson:"id"`
	Data        *string `json:"data" bson:"data"`
	ReceivedAt  string  `json:"received_at" bson:"received_at"`
	ProcessedAt *string `json:"processed_at" bson:"processed_at"`
}

type OutgoingMessageSkriningTB struct {
	ID        string `json:"id" bson:"id"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
	// CkgID     string `json:"ckg_id" bson:"ckg_id"`
}
