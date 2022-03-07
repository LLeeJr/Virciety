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
import {NOTIFY_CONTACT_PERSONS, NOTIFY_HOST_OF_EVENT} from "../service/gql-request-strings";
import {MatSnackBar} from "@angular/material/snack-bar";

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
  selectedEvents: Event[] = [];

  username: string;
  selectedList: string = 'Upcoming events';
  lists: string[] = ['Upcoming events', 'Ongoing events', 'Past events'];

  asc = (a: Event, b: Event) => +new Date(a.startDate) - +new Date(b.startDate);
  desc = (a: Event, b: Event) => +new Date(a.startDate) + +new Date(b.startDate)

  constructor(private keycloak: KeycloakService,
              private gqlService: GQLService,
              private dialog: MatDialog,
              private snackbar: MatSnackBar) { }

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
            this.selectedEvents = this.upcomingEvents;
          }, (error: any) => {
            console.error('there was an error sending the getEvents-query', error);
          });
        })
      } else {
        this.keycloak.login();
      }
    });

    /*this.gqlService.getEvents('user3').subscribe(({data}: any) => {
      // console.log(data);
      this.ongoingEvents = EventComponent.sortEvents(data.getEvents.ongoingEvents, this.asc);
      this.pastEvents = EventComponent.sortEvents(data.getEvents.pastEvents, this.desc);
      this.upcomingEvents = EventComponent.sortEvents(data.getEvents.upcomingEvents, this.asc);
      this.selectedEvents = this.ongoingEvents;
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

  showPeople(people: string[]) {
    this.dialog.open(DialogSubscribersComponent, {
      data: people
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
    if (event.subscribers.indexOf(this.username) < 0) {
      event.subscribers = [...event.subscribers, this.username];
    } else {
      event.subscribers = event.subscribers.filter(member => member !== this.username);
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
    let left: boolean;
    if (event.attendees.indexOf(this.username) < 0) {
      left = false;
      event.attendees = [...event.attendees, this.username];
    } else {
      left = true;
    }

    this.gqlService.attendEvent(event, left, this.username).subscribe(({data}: any) => {
        if (data.attendEvent === "success") {
          event.currentlyAttended = !event.currentlyAttended;
        }},
      ((error: Error) => {
        if (error.message === 'event is expired') {
          this.snackbar.open(error.message + ', please reload the page', undefined,{
            duration: 5 * 1000,
          })
        } else {
          this.snackbar.open('Sorry! An error occurred.', undefined,{
            duration: 5 * 1000,
          })
        }}
    ));
  }

  reportCovidCase(attendees: string[], id: string) {
    let dialogRef = this.dialog.open(DialogReportCovidCaseComponent, {
      data: attendees
    });

    dialogRef.afterClosed().subscribe((covidCase: string | null) => {
      if (covidCase) {
        this.gqlService.notify(NOTIFY_CONTACT_PERSONS, covidCase, id).subscribe(({data}: any) => {
        });
      }
    });
  }

  notifyHost(host: string, id: string) {
    this.gqlService.notify(NOTIFY_HOST_OF_EVENT, host, id, this.username).subscribe(({data}: any) => {
      console.log(data);
    });
  }
}

@Component({
  selector: 'dialog-subscribers',
  template: `
    <h1 mat-dialog-title>subscribers</h1>
    <div mat-dialog-content>
        <mat-list>
            <mat-list-item role="listitem" *ngFor="let username of data">{{username}}</mat-list-item>
        </mat-list>
    </div>`,
})
export class DialogSubscribersComponent {
  constructor(@Inject(MAT_DIALOG_DATA) public data: string[]) {}
}

@Component({
  selector: 'dialog-report-covid-case',
  template: `
    <h1 mat-dialog-title>Report covid case</h1>
    <h4>All contact persons will automatically be notified</h4>
    <mat-dialog-content>
      <mat-form-field>
        <mat-label>Attendees</mat-label>
        <mat-select [(value)]="covidCase">
          <mat-option>None</mat-option>
          <mat-option *ngFor="let attendee of data" [value]="attendee">
            {{attendee}}
          </mat-option>
        </mat-select>
      </mat-form-field>
    </mat-dialog-content>

    <mat-dialog-actions class="mat-dialog-actions">
      <button mat-button [mat-dialog-close]="null">Cancel</button>
      <button mat-button [disabled]="covidCase === undefined" [mat-dialog-close]="covidCase">Report</button>
    </mat-dialog-actions>
  `,
  styles: [`
  .mat-dialog-actions {
  justify-content: flex-end;}`]
})
export class DialogReportCovidCaseComponent {
  covidCase: string;

  constructor(@Inject(MAT_DIALOG_DATA) public data: string[]) {}
}
