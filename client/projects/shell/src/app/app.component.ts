import {Component, OnInit} from '@angular/core';
import {KeycloakService} from "keycloak-angular";
import {AuthLibService} from "auth-lib";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";
import {Router} from "@angular/router";
import {MatSnackBar} from "@angular/material/snack-bar";
import {environment} from "../environments/environment";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'shell';
  isLoggedIn: boolean = false;
  username: string

  isPhonePortrait: boolean = false;
  userMfeOnline: boolean = true;
  notificationMfeOnline: boolean = true;
  postMfeOnline: boolean = true;
  private durationTime: number = 3;

  postMFE: string;
  notifsMFE: string;
  userMFE: string;

  constructor(
    private auth: AuthLibService,
    private keycloak: KeycloakService,
    private responsive: BreakpointObserver,
    private router: Router,
    private snackbar: MatSnackBar,
  ) {
    this.postMFE = environment.postMFE;
    this.notifsMFE = environment.notifsMFE;
    this.userMFE = environment.userMFE;
  }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      this.isLoggedIn = loggedIn;
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
        })
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

  handleError(event: any) {
    let {error, component} = event;
    if (error) {
      let msg = `${component} is currently offline!`;
      switch (component) {
        case 'post':
          this.postMfeOnline = false;
          this.router.navigate(['/page-not-found', msg]);
          break;
        case 'notification':
          this.notificationMfeOnline = false;
          break;
        case 'user':
          this.userMfeOnline = false;
          break;
      }
    }
  }

  placeholderHandler(component: string) {
    let msg = `${component} is currently offline!`;
    this.snackbar.open(msg, undefined, {duration: this.durationTime*1000});
  }
}
