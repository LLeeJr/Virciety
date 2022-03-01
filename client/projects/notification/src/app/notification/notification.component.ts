import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";
import {KeycloakService} from "keycloak-angular";

@Component({
  selector: 'app-notification',
  templateUrl: './notification.component.html',
  styleUrls: ['./notification.component.scss'],
  exportAs: 'NotificationComponent',
})
export class NotificationComponent implements OnInit {

  private username: string;
  show = false;
  notifications: {
    event: string,
    id: string,
    params: {
      key: string,
      value: string,
    }[],
    read: boolean,
    receiver: string,
    route: string,
    text: string,
  }[] = [];

  constructor(private api: ApiService,
              private keycloak: KeycloakService) { }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
          this.api.getNotifs(this.username).subscribe((value: any) => {
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
}
