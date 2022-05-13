package server

const TinkUri = "https://api.tink.com"

const TinkTokenUri = "https://api.tink.com/api/v1/oauth/token"

const TinkAdminRoles = "authorization:grant,user:create"

const LinkBankAccountUriFormat = "https://link.tink.com/1.0/transactions/connect-accounts?client_id=%s&redirect_uri=%s&market=%s&locale=%s&authorization_code=%s"

const GetAuthorizeGrantDelegateCodeRoles = "authorization:read,authorization:grant,credentials:refresh,credentials:read,credentials:write,providers:read,user:read"

const GetAuthorizeGrantCodeRoles = "transactions:read,accounts:read,provider-consents:read,user:read"
