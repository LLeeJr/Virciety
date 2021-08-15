import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {

  constructor(private api: ApiService) {
  }

  ngOnInit(): void {
    this.api.getDms().subscribe(value => console.log(value));
  }

}
