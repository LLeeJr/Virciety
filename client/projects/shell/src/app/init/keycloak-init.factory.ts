import {KeycloakService} from "keycloak-angular";


export function initializeKeycloak(
  keycloak: KeycloakService
) {
  return () =>
    keycloak.init({
      config: {
        url: 'http://localhost:8080' + '/auth',
        realm: 'virciety',
        clientId: 'virciety-frontend',
      },
      initOptions: {
        checkLoginIframe: true,
        checkLoginIframeInterval: 60
      }
    });
}
