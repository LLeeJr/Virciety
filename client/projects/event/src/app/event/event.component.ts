import { Component, OnInit } from '@angular/core';
import {MatDatepickerInputEvent} from "@angular/material/datepicker";
import {formatDate} from "@angular/common";
import {KeycloakService} from "keycloak-angular";
import {GQLService} from "../service/gql.service";

@Component({
  selector: 'app-event',
  templateUrl: './event.component.html',
  styleUrls: ['./event.component.scss']
})
export class EventComponent implements OnInit {

  events: any[] = [];

  startDate: string;
  endDate: string;
  description: string;
  location: string;
  username: string;

  constructor(private keycloak: KeycloakService,
              private gqlService: GQLService) { }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
          this.gqlService.getEvents().subscribe(({data}: any) => {
            console.log(data);
            this.events = data.getEvents;
          }, (error: any) => {
            console.error('there was an error sending the getEvents-query', error);
          });
        })
      } else {
        this.keycloak.login();
      }
    });
  }

  addEvent(type: string, event: MatDatepickerInputEvent<unknown, unknown | null>) {
    if (type === 'start date') {
      this.startDate = formatDate(event.value as string, 'fullDate', 'en-GB');
      //this.events.push(`${type}: ${this.startDate}`);
    } else if (type === 'end date') {
      this.endDate = formatDate(event.value as string, 'fullDate', 'en-GB');
      //this.events.push(`${type}: ${this.endDate}`);
    }
  }

  createEvent() {
    this.gqlService.createEvent(this.username, this.startDate, this.endDate, this.location, this.description);
  }
}
