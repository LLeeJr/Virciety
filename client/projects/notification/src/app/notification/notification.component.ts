import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";
import {KeycloakService} from "keycloak-angular";
import {DatePipe} from "@angular/common";
import {Notification} from "../data/notification";
import {Router} from "@angular/router";

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
            if (value) {
              console.log(value);
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
    switch (n.route) {
      case '/chat':
        let room = {
          name: n.params[1].value,
          id: n.params[2].value,
        }
        sessionStorage.setItem("room", JSON.stringify(room));
        this.router.navigate([`${n.route}/${n.params[1].value}`]).then(() => this.show = !this.show);
        break;
      default:
        break;
    }
  }
}
