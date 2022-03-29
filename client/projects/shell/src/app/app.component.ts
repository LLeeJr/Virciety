import {Component, OnInit} from '@angular/core';
import {KeycloakService} from "keycloak-angular";
import {AuthLibService} from "auth-lib";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";
import {Router} from "@angular/router";
import {MatSnackBar} from "@angular/material/snack-bar";

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
    private responsive: BreakpointObserver,
    private router: Router,
  ) { }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
        })
      } else {
        this.keycloak.login();
      }
    });

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
    let {error, component} = event;
    if (error) {
      let msg = `${component} is currently offline!`;
      switch (component) {
        case 'post':
          this.router.navigate(['/page-not-found', msg]);
      }
    }
  }
}
