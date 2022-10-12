const DOMAIN = 'qa.api.shatteredrealmsonline.com';
const ACCOUNTS_API_VERSION = 'v1';
const ACCOUNTS_API_BASE_URL = `https://${DOMAIN}/accounts`;
const ACCOUNTS_API_URL = `${ACCOUNTS_API_BASE_URL}/${ACCOUNTS_API_VERSION}`

export const environment = {
  production: true,
  DOMAIN: DOMAIN,
  ACCOUNT_API_VERSION: ACCOUNTS_API_VERSION,
  ACCOUNT_API_BASE_URL: ACCOUNTS_API_BASE_URL,
  ACCOUNT_API_URL: ACCOUNTS_API_URL,
};
