export const DEFAULT_HEADERS = {
  'accept': 'application/json',
  'content-type': 'application/json'
}

export const CORE_URI = process.env.CORE_URI || "http://localhost:8081/api"

export const FAVA_URI = process.env.FAVA_URI || "http://localhost:5000"

export const LOGIN_PATH = '/login'

export const REGISTER_PATH = '/register'

export const SYNC_PATH = '/sync-ledger'

export const LINK_PATH = '/auth/link'

export const SESS_ID = 'sessionId'

export const USERNAME = 'username'

export const TOKEN_PATH = "/auth/token"

export const ACCESS_TOKEN = 'access_token'


