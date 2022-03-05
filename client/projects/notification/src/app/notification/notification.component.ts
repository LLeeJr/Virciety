import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";
import {KeycloakService} from "keycloak-angular";
import {DatePipe} from "@angular/common";
import {Notification} from "../data/notification";
import {Router} from "@angular/router";
import {take} from "rxjs";

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

  constructor(private api: ApiService,
              private datePipe: DatePipe,
              private keycloak: KeycloakService,
              private router: Router) { }

  async ngOnInit(): Promise<void> {
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
              if (this.notifications.length == 10) {
                let slice = this.notifications.slice(0, this.notifications.length-1);
                this.notifications = [...slice];
              }
              this.notifications = [value.data.notifAdded, ...this.notifications];
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
      case '/event':
        this.router.navigate([n.route], { queryParams: {eventId: n.params[0].value}}).then(() => this.show = !this.show);
        break;
      default:
        break;
    }
  }

  countUnread() {
    return this.notifications.filter(n => !n.read).length;
  }
}
