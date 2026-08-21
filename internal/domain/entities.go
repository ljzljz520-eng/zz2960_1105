package domain

type BatchStatus string

const (
	BatchDraft     BatchStatus = "draft"
	BatchReview    BatchStatus = "review"
	BatchConfirmed BatchStatus = "confirmed"
	BatchPublished BatchStatus = "published"
	BatchArchived  BatchStatus = "archived"
)

type Batch struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Owner       string      `json:"owner"`
	Status      BatchStatus `json:"status"`
	Version     int         `json:"version"`
	CreatedBy   string      `json:"created_by"`
	RecordCount int         `json:"record_count"`
}

type Record struct {
	ID        string `json:"id"`
	BatchID   string `json:"batch_id"`
	Label     string `json:"label"`
	Expected  int    `json:"expected"`
	Observed  int    `json:"observed"`
	Result    string `json:"result"`
	Confirmed bool   `json:"confirmed"`
	Published bool   `json:"published"`
	Version   int    `json:"version"`
	UpdatedBy string `json:"updated_by"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	BatchID  string `json:"batch_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Detail   string `json:"detail"`
	Sequence int    `json:"sequence"`
}

type CollaborationNote struct {
	ID       string `json:"id"`
	BatchID  string `json:"batch_id"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	Resolved bool   `json:"resolved"`
	Sequence int    `json:"sequence"`
}

type ExportSnapshot struct {
	ID          string `json:"id"`
	BatchID     string `json:"batch_id"`
	Digest      string `json:"digest"`
	Payload     string `json:"payload"`
	RecordCount int    `json:"record_count"`
	PublishedBy string `json:"published_by"`
	Sequence    int    `json:"sequence"`
}
