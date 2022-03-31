import { Component, OnInit } from '@angular/core';
import {Router} from "@angular/router";
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-event',
  templateUrl: './event.component.html',
  styleUrls: ['./event.component.scss']
})
export class EventComponent implements OnInit {
  eventMFE: string;

  constructor(private router: Router) {
    this.eventMFE = environment.eventMFE;
  }

  ngOnInit(): void {
  }

  handleError(event: any) {
    let {error, component} = event;
    if (error && component === 'event') {
      let msg = `${component} is currently offline!`;
      this.router.navigate(['/page-not-found', msg])
    }
  }
}
