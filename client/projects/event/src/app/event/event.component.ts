import { Component, OnInit } from '@angular/core';
import {MatDatepickerInputEvent} from "@angular/material/datepicker";
import {formatDate} from "@angular/common";
import {KeycloakService} from "keycloak-angular";
import {GQLService} from "../service/gql.service";
import {Event} from "../model/event";

@Component({
  selector: 'app-event',
  templateUrl: './event.component.html',
  styleUrls: ['./event.component.scss']
})
export class EventComponent implements OnInit {

  events: Event[] = [];

  startDate: string;
  endDate: string;
  description: string;
  location: string;
  username: string;
  title: string;

  constructor(private keycloak: KeycloakService,
              private gqlService: GQLService) { }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
          this.gqlService.getEvents().subscribe(({data}: any) => {
            // console.log(data);
            const unsorted: Event[] = [];

            for (let getEvent of data.getEvents) {
              const event: Event = new Event(getEvent);

              unsorted.push(event);
            }

            this.events = unsorted.sort((a, b) => +new Date(b.startDate) - +new Date(a.startDate))
          }, (error: any) => {
            console.error('there was an error sending the getEvents-query', error);
          });
        })
      } else {
        this.keycloak.login();
      }
    });
  }

  // create dialog for this
  // ideas: click on location and redirect to google maps
  // https://developers.google.com/maps/documentation/urls/get-started#search-action
  createEvent() {
    this.gqlService
      .createEvent(this.title, this.username, this.startDate, this.endDate, this.location, this.description)
      .subscribe(({data}: any) => {
        console.log(data);
        this.events = [...this.events, data.createEvent]
      });
  }

  saveStartDate(event: MatDatepickerInputEvent<unknown, unknown | null>) {
    this.startDate = formatDate(event.value as string, 'fullDate', 'en-GB');
  }

  saveEndDate(event: MatDatepickerInputEvent<unknown, unknown | null>) {
    this.endDate = formatDate(event.value as string, 'fullDate', 'en-GB');
  }
}
