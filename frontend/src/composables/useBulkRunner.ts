import type { BulkFailure } from '@/types/bulk'

export interface BulkRunResult {
  total: number
  success: number
  failures: BulkFailure[]
}

export async function runBulk(items: string[], handler: (item: string) => Promise<void>): Promise<BulkRunResult> {
  const failures: BulkFailure[] = []
  let success = 0

  for (const item of items) {
    try {
      await handler(item)
      success += 1
    } catch (error: any) {
      failures.push({
        item,
        reason: error?.message || 'Unknown error',
      })
    }
  }

  return {
    total: items.length,
    success,
    failures,
  }
}
