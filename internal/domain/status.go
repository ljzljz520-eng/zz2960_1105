package domain

func IsFinal(status BatchStatus) bool {
	return status == BatchPublished || status == BatchArchived
}

func IsEditable(status BatchStatus) bool {
	return status == BatchDraft || status == BatchReview
}

func NextStatus(status BatchStatus) BatchStatus {
	switch status {
	case BatchDraft:
		return BatchReview
	case BatchReview:
		return BatchConfirmed
	case BatchConfirmed:
		return BatchPublished
	case BatchPublished:
		return BatchArchived
	default:
		return status
	}
}
