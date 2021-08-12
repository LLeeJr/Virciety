import {Component, OnInit} from '@angular/core';
import {ApiService} from "./api/api.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'chat';

  constructor(private api: ApiService) {
  }

  ngOnInit(): void {
    this.api.getDms().subscribe(value => console.log(value));
  }
}
