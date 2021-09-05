import {Component, OnInit} from '@angular/core';
import {ApiService} from "./api/api.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'chat';

  constructor() {
  }

}
