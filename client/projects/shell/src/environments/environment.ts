// This file can be replaced during build by using the `fileReplacements` array.
// `ng build` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

export const environment = {
  production: false,
  chatMFE: "http://localhost:5001/remoteEntry.js",
  eventMFE: "http://localhost:5004/remoteEntry.js",
  keycloak: "http://localhost:8080/auth",
  postMFE: "http://localhost:5002/remoteEntry.js",
  profileMFE: "http://localhost:5005/remoteEntry.js",
  notifsMFE: "http://localhost:5006/remoteEntry.js",
  userMFE: "http://localhost:5003/remoteEntry.js",
  userAPI: "http://localhost:8085/query"
};

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/plugins/zone-error';  // Included with Angular CLI.
