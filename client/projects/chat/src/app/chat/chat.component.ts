import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";
import {AuthLibService} from "auth-lib";

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {

  openChats: {
    __typename: string,
    preview: string,
    withUser: string,
  }[] = [];

  constructor(private api: ApiService,
              private auth: AuthLibService) {
  }

  ngOnInit(): void {
    this.auth.getUserName().subscribe(user => {
      if (user !== '') {
        this.api.getOpenChats(user).subscribe(value => {
          this.openChats = value.data.getOpenChats;
        });
      }
    });

    if (this.auth.userName !== '') {
      this.api.getOpenChats(this.auth.userName).subscribe(value => {
        this.openChats = value.data.getOpenChats;
      });
    }
  }

  setChatPartner(withUser: string) {
    this.api.chatPartner = withUser;
  }
}
