import {Component, OnInit} from '@angular/core';
import {KeycloakService} from "keycloak-angular";
import {GQLService} from "../service/gql.service";
import {Event} from "../model/event";
import {MatDialog} from "@angular/material/dialog";
import {CreateEventComponent} from "../create-event/create-event.component";

@Component({
  selector: 'app-event',
  templateUrl: './event.component.html',
  styleUrls: ['./event.component.scss']
})
export class EventComponent implements OnInit {

  events: Event[] = [];

  username: string;

  constructor(private keycloak: KeycloakService,
              private gqlService: GQLService,
              private dialog: MatDialog) { }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
          this.gqlService.getEvents().subscribe(({data}: any) => {
            console.log(data);
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

    /*this.gqlService.getEvents().subscribe(({data}: any) => {
      // console.log(data);
      const unsorted: Event[] = [];

      for (let getEvent of data.getEvents) {
        const event: Event = new Event(getEvent);

        unsorted.push(event);
      }

      this.events = unsorted.sort((a, b) => +new Date(b.startDate) - +new Date(a.startDate))
    }, (error: any) => {
      console.error('there was an error sending the getEvents-query', error);
    });*/
  }

  createEvent() {
    /*let startDate = formatDate(this.range.controls.start.value, 'fullDate', 'en-GB');
    let endDate = formatDate(this.range.controls.end.value, 'fullDate', 'en-GB');
    if (this.checked) {
      let shortDate = formatDate(startDate, 'short', 'en-GB');
      let split = shortDate.split(',', 1);
      startDate = split[0] + ', ' + this.startTime.value;

      shortDate = formatDate(endDate, 'short', 'en-GB');
      split = shortDate.split(',', 1);
      endDate = split[0] + ', ' + this.endTime.value;
    }


    this.gqlService
      .createEvent(this.title.value, this.username, startDate, endDate, this.location, this.description)
      .subscribe(({data}: any) => {
        console.log(data);
        this.events = [...this.events, data.createEvent].sort((a, b) => +new Date(b.startDate) - +new Date(a.startDate))
      });*/
  }

  gotToMaps(location: string) {
    window.open('https://www.google.com/maps/search/?api=1&query=' + encodeURIComponent(location))
  }

  openDialog(event: Event | null, editMode: boolean) {
    let dialogRef = this.dialog.open(CreateEventComponent, {
      data: {
        editMode: editMode,
        event: event,
      }
    })

    dialogRef.afterClosed().subscribe(data => {
      console.log(data);
    })
  }
}
