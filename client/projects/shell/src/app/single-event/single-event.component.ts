import { Component, OnInit } from '@angular/core';
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-single-event',
  templateUrl: './single-event.component.html',
  styleUrls: ['./single-event.component.scss']
})
export class SingleEventComponent implements OnInit {
  eventMFE: string;

  constructor() { }

  ngOnInit(): void {
    this.eventMFE = environment.eventMFE;
  }

}
