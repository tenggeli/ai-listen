export interface UserOrderFeedback {
  id: string
  orderId: string
  userId: string
  ratingScore: number
  reviewTags: string[]
  reviewContent: string
  hasComplaint: boolean
  complaintReason: string
  complaintContent: string
  createdAt: string
}
