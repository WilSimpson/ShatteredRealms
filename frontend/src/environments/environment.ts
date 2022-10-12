// This file can be replaced during build by using the `fileReplacements` array.
// `ng build` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.
const DOMAIN = 'localhost';
const ACCOUNT_API_PORT = 8080;
const ACCOUNT_API_VERSION = 'v1';
const ACCOUNT_API_BASE_URL = `http://${DOMAIN}:${ACCOUNT_API_PORT}`;
const ACCOUNT_API_URL = `${ACCOUNT_API_BASE_URL}/${ACCOUNT_API_VERSION}`

export const environment = {
  production: false,
  DOMAIN: DOMAIN,
  ACCOUNT_API_VERSION: ACCOUNT_API_VERSION,
  ACCOUNT_API_BASE_URL: ACCOUNT_API_BASE_URL,
  ACCOUNT_API_URL: ACCOUNT_API_URL,
  recaptcha: {
    siteKey: '6LeFKkEhAAAAALeg05zJHu1kKaj-vSPdJ80ilEBV'
  },
};



/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/plugins/zone-error';  // Included with Angular CLI.
