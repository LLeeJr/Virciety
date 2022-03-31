import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";
import {KeycloakService} from "keycloak-angular";
import {DatePipe} from "@angular/common";
import {Notification} from "../data/notification";
import {Router} from "@angular/router";
import {take} from "rxjs";
import {MatSnackBar} from "@angular/material/snack-bar";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";

@Component({
  selector: 'app-notification',
  templateUrl: './notification.component.html',
  styleUrls: ['./notification.component.scss'],
  exportAs: 'NotificationComponent',
})
export class NotificationComponent implements OnInit {

  private username: string;
  show = false;
  notifications: Notification[] = [];
  isPhonePortrait: boolean = false;

  constructor(private api: ApiService,
              private datePipe: DatePipe,
              private keycloak: KeycloakService,
              private router: Router,
              private snackbar: MatSnackBar,
              private responsive: BreakpointObserver) { }

  async ngOnInit(): Promise<void> {
    this.responsive.observe(Breakpoints.HandsetPortrait).subscribe((result) => {
      this.isPhonePortrait = result.matches;
    });

    this.api.errorState.subscribe(value => this.snackbar.open(value, undefined, {duration: 3000}));
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
          this.api.getNotifications(this.username).subscribe((value: any) => {
            if (value && value.data && value.data.getNotifsByReceiver) {
              this.notifications = value.data.getNotifsByReceiver;
            }
          });

          this.api.subscribeToNotifications(this.username).subscribe((value: any) => {
            if (value && value.data && value.data.notifAdded) {
              let notification = value.data.notifAdded;
              if (this.notifications.length == 10) {
                let slice = this.notifications.slice(0, this.notifications.length-1);
                this.notifications = [...slice];
              }
              this.notifications = [notification, ...this.notifications];

              let url = this.router.url;
              let routes = url.split('/')

              if (notification.route === "/chat" && url.startsWith(notification.route)) {
                if (notification.params[1] && notification.params[1].key === "roomName" && routes.includes(notification.params[1].value)) {
                  this.api.setReadStatus(notification.id, true).pipe(take(1)).subscribe(() => {
                    notification.read = true;
                  });
                }
              }
            }
          });
        })
      } else {
        this.keycloak.login();
      }
    })
  }

  showNotifications() {
    this.show = !this.show;
  }

  transformDate(date: string) {
    return this.datePipe.transform(date, 'short')
  }

  clickNotification(n: Notification) {
    if (!n.read) {
      this.api.setReadStatus(n.id, !n.read).pipe(take(1)).subscribe(value => {
        if (value && value.data && value.data.setReadStatus) {
          let {id, read} = value.data.setReadStatus;
          let arr = [];
          for (let n of this.notifications) {
            if (n.id === id) {
              n = { ...n, read: read};
            }
            arr.push(n);
          }
          this.notifications = [...arr];
        }
      });
    }
    switch (n.route) {
      case '/chat':
        let room = {
          name: n.params[1].value,
          id: n.params[2].value,
        }
        sessionStorage.setItem("room", JSON.stringify(room));
        this.router.navigate([`${n.route}/${n.params[1].value}`]).then(() => this.show = !this.show);
        break;
      case '/profile':
        let user = n.params[0].key == "newFollower" ? n.params[0].value : n.receiver;
        this.router.navigate([n.route], { queryParams: {username: user}}).then(() => this.show = !this.show);
        break;
      case '/e/':
        this.router.navigateByUrl('/', {skipLocationChange: true}).then(() =>
          this.router.navigate([n.route+n.params[0].value]).then(() =>
            this.show = !this.show
          )
        );
        break;
      case '/p/':
        let postIDIndex = n.params.length > 1 ? 1 : 0;
        this.router.navigateByUrl('/', {skipLocationChange: true}).then(() =>
          this.router.navigate([n.route+n.params[postIDIndex].value]).then(() =>
            this.show = !this.show
          )
        );
        break;
      default:
        break;
    }
  }

  countUnread() {
    return this.notifications.filter(n => !n.read).length;
  }
}
