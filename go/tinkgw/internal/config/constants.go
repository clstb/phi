package config

const TinkApiUri = "https://api.tink.com"

const TinkAdminRoles = "authorization:grant" //#,user:create"

const DelegatedAuthorizationRoles = "authorization:read,authorization:grant,credentials:refresh,credentials:read,credentials:write,providers:read,user:read"

const AuthorizeGrantRoles = "transactions:read,accounts:read,provider-consents:read,user:read"

const AuthorizationCodeGrantType = "authorization_code"

const JsonMediaType = "application/json"

const DefaultMarket = "DE"

const DefaultLocale = "de_DE"

const DelegatedAuthorizationPath = "/api/v1/oauth/authorization-grant/delegate"

const UserCreatePath = "/api/v1/user/create"

const TransactionsPath = "/data/v2/transactions"

const AccountsPath = "/data/v2/accounts"

const ProvidersPath = "/api/v1/providers"

const AccessTokenPath = "/api/v1/oauth/token"
