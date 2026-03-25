package complaint

type CreateComplaintRequest struct {
	ComplaintType  string   `json:"complaintType"`
	Content        string   `json:"content"`
	EvidenceImages []string `json:"evidenceImages"`
}
