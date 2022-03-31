import {KeycloakService} from "keycloak-angular";
import {environment} from "../../environments/environment";


export function initializeKeycloak(
  keycloak: KeycloakService
) {
  return () =>
    keycloak.init({
      config: {
        url: environment.keycloak,
        realm: 'virciety',
        clientId: 'virciety-frontend',
      },
      initOptions: {
        // checkLoginIframe: true,
        // checkLoginIframeInterval: 60
      }
    });
}
