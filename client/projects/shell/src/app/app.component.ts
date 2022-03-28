import {Component, OnInit} from '@angular/core';
import {KeycloakService} from "keycloak-angular";
import {AuthLibService} from "auth-lib";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'shell';
  isLoggedIn: boolean = false;
  username: string

  searchMode: boolean = false;
  isPhonePortrait: boolean = false;

  constructor(
    private auth: AuthLibService,
    private keycloak: KeycloakService,
    private responsive: BreakpointObserver
  ) {
    this.keycloak.isLoggedIn().then((loggedIn) => {
      this.isLoggedIn = loggedIn;
      this.keycloak.loadUserProfile().then(() => this.username = this.keycloak.getUsername())
    });
  }

  ngOnInit() {
    this.responsive.observe(Breakpoints.HandsetPortrait).subscribe((result) => {
      this.isPhonePortrait = result.matches;
    });
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

  openSearch() {
    this.searchMode = !this.searchMode;
  }

  closeSearch() {
    this.searchMode = !this.searchMode;
  }

  handleError(event: any) {
    console.log(event)
  }
}
