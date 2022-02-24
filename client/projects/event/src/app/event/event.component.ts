import {Component, Inject, OnInit} from '@angular/core';
import {KeycloakService} from "keycloak-angular";
import {GQLService} from "../service/gql.service";
import {Event} from "../model/event";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";
import {CreateEventComponent} from "../create-event/create-event.component";
import {CreateEventData} from "../model/dialog.data";
import {formatDate} from "@angular/common";
import {ContactDetailsComponent} from "../contact-details/contact-details.component";
import {UserData} from "../model/userData";

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

  upcomingEvents: Event[];
  ongoingEvents: Event[];
  pastEvents: Event[];
  selectedEvents: Event[];

  username: string;
  selectedList: string = 'Ongoing events'; // TODO
  lists: string[] = ['Upcoming events', 'Ongoing events', 'Past events'];

  asc = (a: Event, b: Event) => +new Date(a.startDate) - +new Date(b.startDate);
  desc = (a: Event, b: Event) => +new Date(a.startDate) + +new Date(b.startDate)

  constructor(private keycloak: KeycloakService,
              private gqlService: GQLService,
              private dialog: MatDialog) { }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
          this.gqlService.getEvents(this.username).subscribe(({data}: any) => {
            // console.log(data);
            this.ongoingEvents = EventComponent.sortEvents(data.getEvents.ongoingEvents, this.asc);
            this.pastEvents = EventComponent.sortEvents(data.getEvents.pastEvents, this.desc);
            this.upcomingEvents = EventComponent.sortEvents(data.getEvents.upcomingEvents, this.asc);
            this.selectedEvents = this.ongoingEvents; // TODO
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
      this.ongoingEvents = EventComponent.sortEvents(data.getEvents.ongoingEvents, this.asc);
      this.pastEvents = EventComponent.sortEvents(data.getEvents.pastEvents, this.desc);
      this.upcomingEvents = EventComponent.sortEvents(data.getEvents.upcomingEvents, this.asc);
      this.selectedEvents = this.ongoingEvents; // TODO
    }, (error: any) => {
      console.error('there was an error sending the getEvents-query', error);
    });*/
  }

  private static sortEvents(events: any, sortBy: (a: Event, b: Event) => number): Event[] {
    const unsorted: Event[] = [];

    for (let getEvent of events) {
      const event: Event = new Event(getEvent);
      unsorted.push(event);
    }

    return unsorted.sort(sortBy);
  }

  gotToMaps(location: string) {
    window.open('https://www.google.com/maps/search/?api=1&query=' + encodeURIComponent(location));
  }

  openEventDialog(event: Event | null, editMode: boolean) {
    let dialogRef = this.dialog.open(CreateEventComponent, {
      data: {
        editMode: editMode,
        event: event,
      }
    });

    dialogRef.afterClosed().subscribe((data: CreateEventData | false) => {
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

  showPeople(members: string[]) {
    this.dialog.open(DialogMembersComponent, {
      data: members
    });
  }

  changeList() {
    if (this.selectedList === 'Upcoming events') {
      this.selectedEvents = this.upcomingEvents;
    } else if (this.selectedList === 'Ongoing events') {
      this.selectedEvents = this.ongoingEvents;
    } else {
      this.selectedEvents = this.pastEvents;
    }
  }

  createEvent(event: Event | { description: string; location: string; startDate: string; endDate: string; startTime: string | null; endTime: string | null; title: string } | null) {
    // console.log(event)

    if (event && !(event instanceof Event)) {
      const eventDate = EventComponent.formatEventDate(event);

      this.gqlService
        .createEvent(event.title, this.username, eventDate.startDate, eventDate.endDate, event.location, event.description)
        .subscribe(({data}: any) => {
          // console.log(data);

          if (data.createEvent.type === 'upcoming') {
            this.upcomingEvents = [...this.upcomingEvents, new Event(data.createEvent.event)].sort(this.asc);
          } else if (data.createEvent.type === 'ongoing') {
            this.ongoingEvents = [...this.ongoingEvents, new Event(data.createEvent.event)].sort(this.asc);
          }

          if (this.selectedList.toLowerCase().includes(data.createEvent.type)) {
            this.selectedEvents = [...this.selectedEvents, new Event(data.createEvent.event)].sort(this.asc);
          }
        });
    } else {
      console.error('CreateEvent \'event\' is instance of event or null')
    }
  }

  removeEvent(event: Event | { description: string; location: string; startDate: string; endDate: string; startTime: string | null; endTime: string | null; title: string } | null) {
    if (event && event instanceof Event) {
      this.gqlService.removeEvent(event.id).subscribe(({data}: any) => {
        if (data.removeEvent === "success") {
          if (this.selectedList === 'Upcoming events') {
            this.upcomingEvents = this.upcomingEvents.filter(e => e.id !== event.id);
            this.selectedEvents = this.upcomingEvents;
          } else if (this.selectedList === 'Ongoing events') {
            this.ongoingEvents = this.ongoingEvents.filter(e => e.id !== event.id);
            this.selectedEvents = this.ongoingEvents;
          } else {
            this.pastEvents = this.pastEvents.filter(e => e.id !== event.id);
            this.selectedEvents = this.pastEvents;
          }
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
        // console.log(data);
        if (data.editEvent.type === 'ongoing' && this.selectedList === 'Upcoming events') {
          this.ongoingEvents = [...this.ongoingEvents, event].sort(this.asc);
          this.upcomingEvents = this.upcomingEvents.filter(upcomingEvent => upcomingEvent.id !== event.id);
          this.selectedEvents = this.upcomingEvents;
        } else if (data.editEvent.type === 'upcoming' && this.selectedList === 'Ongoing events') {
          this.upcomingEvents = [...this.upcomingEvents, event].sort(this.asc);
          this.ongoingEvents = this.ongoingEvents.filter(ongoingEvent => ongoingEvent.id !== event.id);
          this.selectedEvents = this.ongoingEvents;
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

  subscribeEvent(event: Event) {
    if (event.members.indexOf(this.username) < 0) {
      event.members = [...event.members, this.username];
    } else {
      event.members = event.members.filter(member => member !== this.username);
    }

    this.gqlService.subscribeEvent(event).subscribe(({data}: any) => {
      // console.log(data);
    });
  }

  checkUserData(event: Event) {
    this.gqlService.userDataExists(this.username).subscribe(({data}: any) => {
      // console.log(data);
      if (data.userDataExists) {
        this.attendEvent(event);
      } else {
        let dialogRef = this.dialog.open(ContactDetailsComponent, {
          data: this.username
        });

        dialogRef.afterClosed().subscribe((userData: UserData | null) => {
          if (!userData) {
            return;
          }

          this.gqlService.addUserData(userData).subscribe(({data}: any) => {
            this.attendEvent(event);
          });
        });
      }
    })
  }

  attendEvent(event: Event) {
    let left;
    if (event.attending.indexOf(this.username) < 0) {
      left = false;
      event.attending = [...event.attending, this.username];
    } else {
      left = true;
    }

    this.gqlService.attendEvent(event, left, this.username).subscribe(({data}: any) => {
      console.log(data);
    })
  }

  reportCovidCase(attending: string[]) {

  }

  notifyHost(attending: string[]) {

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
