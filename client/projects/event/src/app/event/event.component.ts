import {Component, Inject, OnInit} from '@angular/core';
import {KeycloakService} from "keycloak-angular";
import {GQLService} from "../service/gql.service";
import {Event} from "../model/event";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";
import {CreateEventComponent, OutputData} from "../create-event/create-event.component";
import {formatDate} from "@angular/common";

export interface EventDate {
  startDate: string;
  endDate: string;
}

@Component({
  selector: 'app-event',
  templateUrl: './event.component.html',
  styleUrls: ['./event.component.scss']
})
export class EventComponent implements OnInit {

  data: any;
  upcomingEvents: Event[];
  ongoingEvents: Event[];
  pastEvents: Event[];
  selectedEvents: Event[];

  username: string;
  selectedList: string = 'Upcoming events';
  lists: string[] = ['Upcoming events', 'Ongoing events', 'Past events'];

  constructor(private keycloak: KeycloakService,
              private gqlService: GQLService,
              private dialog: MatDialog) { }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
          this.gqlService.getEvents().subscribe(({data}: any) => {
            this.data = data;
            console.log(data);
            this.upcomingEvents = this.sortEvents(data.getEvents.upcomingEvents);
            this.selectedEvents = this.upcomingEvents;
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

  private sortEvents(events: any): Event[] {
    const unsorted: Event[] = [];

    for (let getEvent of events) {
      const event: Event = new Event(getEvent);
      unsorted.push(event);
    }

    return unsorted.sort((a, b) => +new Date(a.startDate) - +new Date(b.startDate));
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

    dialogRef.afterClosed().subscribe((data: OutputData | false) => {
      // console.log(data);
      if (data === undefined || data === false) {
        return;
      }

      // remove event
      if (data.remove && editMode) {
        this.removeEvent(data.event);
      // edit event
      } else if (!data.remove && editMode) {
        this.editEvent(data.event);
      // create event
      } else if (!data.remove && !editMode) {
        this.createEvent(data.event);
      }
    });
  }

  showMembers() {
    this.dialog.open(DialogMembersComponent);
  }

  radioChange() {
    if (this.selectedList === 'Upcoming events') {
      this.selectedEvents = this.upcomingEvents;
    } else if (this.selectedList === 'Ongoing events') {
      if (!this.ongoingEvents) {
        this.ongoingEvents = this.sortEvents(this.data.getEvents.ongoingEvents);
      }
      this.selectedEvents = this.ongoingEvents;
    } else {
      if (!this.pastEvents) {
        this.pastEvents = this.sortEvents(this.data.getEvents.pastEvents);
      }
      this.selectedEvents = this.pastEvents;
    }
  }

  createEvent(event: Event | { description: string; location: string; startDate: string; endDate: string; startTime: string | null; endTime: string | null; title: string } | null) {
    console.log(event)

    if (event && !(event instanceof Event)) {
      const eventDate = EventComponent.formatEventDate(event);

      this.gqlService
        .createEvent(event.title, this.username, eventDate.startDate, eventDate.endDate, event.location, event.description)
        .subscribe(({data}: any) => {
          // console.log(data);
          this.upcomingEvents = [...this.upcomingEvents, new Event(data.createEvent)].sort((a, b) => +new Date(a.startDate) - +new Date(b.startDate))
        });
    } else {
      console.error('CreateEvent \'event\' is instance of event or null')
    }
  }

  private removeEvent(event: Event | { description: string; location: string; startDate: string; endDate: string; startTime: string | null; endTime: string | null; title: string } | null) {
    if (event && event instanceof Event) {
      this.gqlService.removeEvent(event.id).subscribe(({data}: any) => {
        if (data.removeEvent === "success") {
          this.upcomingEvents = this.upcomingEvents.filter(e => e.id !== event.id);
        }
      })
    } else {
      console.error('RemoveEvent \'event\' is not instance of event or null');
    }
  }

  private editEvent(event: Event | { description: string; location: string; startDate: string; endDate: string; startTime: string | null; endTime: string | null; title: string } | null) {
    if (event && event instanceof Event) {
      const eventDate = EventComponent.formatEventDate(event);
      this.gqlService.editEvent(event, eventDate).subscribe(({data}: any) => {
        if (data.removeEvent === "success") {
          this.upcomingEvents = this.upcomingEvents.filter(e => e.id !== event.id);
        }
      })
    } else {
      console.error('EditEvent \'event\' is not instance of event or null')
    }
  }

  private static formatEventDate(event: Event | { description: string; location: string; startDate: string; endDate: string; startTime: string | null; endTime: string | null; title: string }): EventDate {
    let startDate = formatDate(event.startDate, 'fullDate', 'en-GB');
    let endDate = formatDate(event.endDate, 'fullDate', 'en-GB');
    if (event.startTime && event.endTime) {
      let shortDate = formatDate(startDate, 'short', 'en-GB');
      let split = shortDate.split(',', 1);
      startDate = split[0] + ', ' + event.startTime;

      shortDate = formatDate(endDate, 'short', 'en-GB');
      split = shortDate.split(',', 1);
      endDate = split[0] + ', ' + event.endTime;
    }

    return {startDate: startDate, endDate: endDate}
  }
}

@Component({
  selector: 'dialog-members',
  template: `
    <h1 mat-dialog-title>Members</h1>
    <div mat-dialog-content>
        <mat-list>
            <mat-list-item role="listitem" *ngFor="let username of data">{{username}}</mat-list-item>
        </mat-list>
    </div>`,
})
export class DialogMembersComponent {
  constructor(@Inject(MAT_DIALOG_DATA) public data: string[]) {}
}
