import {Component, Inject, OnInit} from '@angular/core';
import {KeycloakService} from "keycloak-angular";
import {GQLService} from "../service/gql.service";
import {Event} from "../model/event";
import {MAT_DIALOG_DATA, MatDialog} from "@angular/material/dialog";
import {CreateEventComponent, OutputData} from "../create-event/create-event.component";

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

            this.events = unsorted.sort((a, b) => +new Date(a.startDate) - +new Date(b.startDate))
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

    dialogRef.afterClosed().subscribe((data: OutputData) => {
      console.log(data);
      if (data === undefined) {
        return;
      }

      // remove event
      if (data.remove && editMode) {
        this.removeEvent(data.event);
      } else if (!data.remove && editMode) {

      }
    })
  }

  showMembers() {
    this.dialog.open(DialogMembersComponent)
  }

  removeEvent(event: Event | { description: string; location: string; startDate: string; endDate: string; startTime: string; endTime: string; title: string } | null) {
    if (event instanceof Event) {
      this.gqlService.removeEvent(event.id).subscribe(({data}: any) => {
        if (data.removeEvent === "success") {
          this.events = this.events.filter(e => e.id !== event.id);
        }
      })
    } else {
      console.log(event, 'is not instance of event')
    }
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
