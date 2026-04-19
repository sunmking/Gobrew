export type Domain = 'dashboard' | 'install' | 'upgrade' | 'maintain' | 'services'

export interface DomainConfig {
  entryPath: string
  paths: readonly string[]
}

export interface DomainDisplay {
  labelKey: string
  titleKey: string
  descriptionKey: string
}

export const DOMAIN_DISPLAY: Record<Domain, DomainDisplay> = {
  dashboard: {
    labelKey: 'domains.dashboard.label',
    titleKey: 'domains.dashboard.title',
    descriptionKey: 'domains.dashboard.description',
  },
  install: {
    labelKey: 'domains.install.label',
    titleKey: 'domains.install.title',
    descriptionKey: 'domains.install.description',
  },
  upgrade: {
    labelKey: 'domains.upgrade.label',
    titleKey: 'domains.upgrade.title',
    descriptionKey: 'domains.upgrade.description',
  },
  maintain: {
    labelKey: 'domains.maintain.label',
    titleKey: 'domains.maintain.title',
    descriptionKey: 'domains.maintain.description',
  },
  services: {
    labelKey: 'domains.services.label',
    titleKey: 'domains.services.title',
    descriptionKey: 'domains.services.description',
  },
}

export const DOMAIN_CONFIG: Record<Domain, DomainConfig> = {
  dashboard: {
    entryPath: '/',
    paths: ['/'],
  },
  install: {
    entryPath: '/install',
    paths: ['/install', '/installed', '/explore'],
  },
  upgrade: {
    entryPath: '/upgrade',
    paths: ['/upgrade', '/update'],
  },
  maintain: {
    entryPath: '/maintain',
    paths: ['/maintain', '/taps', '/cleanup', '/bundle'],
  },
  services: {
    entryPath: '/services',
    paths: ['/services'],
  },
}

export const DOMAIN_ENTRY_PATHS: Record<Domain, string> = {
  dashboard: DOMAIN_CONFIG.dashboard.entryPath,
  install: DOMAIN_CONFIG.install.entryPath,
  upgrade: DOMAIN_CONFIG.upgrade.entryPath,
  maintain: DOMAIN_CONFIG.maintain.entryPath,
  services: DOMAIN_CONFIG.services.entryPath,
}

export const DOMAIN_ORDER: readonly Domain[] = ['dashboard', 'install', 'upgrade', 'maintain', 'services']

const PATH_TO_DOMAIN = new Map<string, Domain>()

for (const [domain, domainConfig] of Object.entries(DOMAIN_CONFIG) as Array<[Domain, DomainConfig]>) {
  for (const path of domainConfig.paths) {
    PATH_TO_DOMAIN.set(path, domain)
  }
}

export function normalizePath(path: string): string {
  const raw = path.trim()
  if (!raw) {
    return '/'
  }

  let normalized = raw

  try {
    if (/^https?:\/\//i.test(raw)) {
      const url = new URL(raw)
      normalized = url.hash.startsWith('#/')
        ? url.hash.slice(1)
        : url.pathname
    }
  } catch {
    // Ignore malformed URLs and normalize as a plain path.
  }

  const hashPathIndex = normalized.indexOf('#/')
  if (hashPathIndex >= 0) {
    normalized = normalized.slice(hashPathIndex + 1)
  }

  const [withoutQuery] = normalized.split('?')
  const [withoutHash] = withoutQuery.split('#')

  normalized = withoutHash || '/'
  if (!normalized.startsWith('/')) {
    normalized = `/${normalized}`
  }

  normalized = normalized.replace(/\/{2,}/g, '/')
  if (normalized.length > 1) {
    normalized = normalized.replace(/\/+$/, '')
  }

  return normalized.toLowerCase()
}

export function resolveDomainFromPath(path: string): Domain {
  return PATH_TO_DOMAIN.get(normalizePath(path)) ?? 'dashboard'
}

export function isDomainPath(path: string): boolean {
  return PATH_TO_DOMAIN.has(normalizePath(path))
}

export function getDomainEntryPath(domain: Domain): string {
  return DOMAIN_ENTRY_PATHS[domain]
}
