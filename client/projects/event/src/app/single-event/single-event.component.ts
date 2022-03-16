import { Component, OnInit } from '@angular/core';
import {GQLService} from "../service/gql.service";
import {AuthLibService} from "auth-lib";
import {Event} from "../model/event";
import {ActivatedRoute} from "@angular/router";
import {Location} from "@angular/common";

@Component({
  selector: 'app-single-event',
  templateUrl: './single-event.component.html',
  styleUrls: ['./single-event.component.scss'],
  exportAs: 'SingleEventComponent',
})
export class SingleEventComponent implements OnInit {

  username: string;
  valid: boolean = true;
  event: Event;

  constructor(private gqlService: GQLService,
              private auth: AuthLibService,
              private route: ActivatedRoute,
              private location: Location) {
    this.username = this.auth.userName;
  }

  ngOnInit(): void {
    let eventID = this.route.snapshot.paramMap.get('id');
    // get postID when opened via dialog
    if (eventID === null) {
      eventID = this.location.path().substring(3);
    }

    if (eventID !== null) {
      this.gqlService.getEventByID(eventID).subscribe({
        next: ({data}: any) => {
          this.event = new Event(data.getEvent);
        },
        error: (_: any) => {
          this.valid = false;
        }
      });
    } else {
      this.valid = false
    }
  }

  gotToMaps(location: string) {
    window.open('https://www.google.com/maps/search/?api=1&query=' + encodeURIComponent(location));
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

}
