import {Component, OnInit} from '@angular/core';
import {KeycloakService} from "keycloak-angular";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'shell';
  isLoggedIn: boolean = false;
  username: string

  constructor(
    private keycloak: KeycloakService,
  ) {
    this.keycloak.isLoggedIn().then((loggedIn) => {
      this.isLoggedIn = loggedIn;
      this.keycloak.loadUserProfile().then(() => this.username = this.keycloak.getUsername())
    });
  }

  ngOnInit() {}

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
