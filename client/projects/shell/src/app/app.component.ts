import {Component, OnInit} from '@angular/core';
import {MfeOptions} from "./mfe/mfe";
import {LookupService} from "./mfe/lookup.service";
import {KeycloakService} from "keycloak-angular";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'shell';
  mfes: MfeOptions[] = [];
  isLoggedIn: boolean = false;
  activeMfe: MfeOptions;

  constructor(
    private lookupService: LookupService,
    private keycloak: KeycloakService
  ) {
    this.keycloak.isLoggedIn().then((loggedIn) => {
      this.isLoggedIn = loggedIn;
    });
  }

  async ngOnInit(): Promise<void> {
    this.mfes = await this.lookupService.lookup();
  }

  activate(mfe: MfeOptions): void {
    this.activeMfe = mfe;
  }

  logout() {
    this.keycloak.isLoggedIn().then((loggedIn) => {
      if (loggedIn) {
        this.keycloak.logout(window.location.origin).then(() => {
          this.isLoggedIn = false;
        });
      }
    })
  }
}
