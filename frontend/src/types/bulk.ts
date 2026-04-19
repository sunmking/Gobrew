export interface BulkFailure {
  item: string
  reason: string
}

export interface BulkSummary {
  action: string
  total: number
  success: number
  failures: BulkFailure[]
  timestamp: number
}
