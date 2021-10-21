import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";
import {AuthLibService} from "auth-lib";
import {Room} from "../data/room";

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {

  chatrooms: Room[] = [];

  constructor(private api: ApiService,
              private auth: AuthLibService) {
  }

  ngOnInit(): void {
    this.auth.getUserName().subscribe(user => {
      if (user !== '') {
        this.api.getRoomsByUser(user).subscribe(value => {
          this.chatrooms = value.data.getRoomsByUser;
        });
      }
    });

    if (this.auth.userName !== '') {
      this.api.getRoomsByUser(this.auth.userName).subscribe(value => {
        this.chatrooms = value.data.getRoomsByUser;
      });
    }
  }

  setChatPartner(room: Room) {
    this.api.chatMembers = room.member;
    this.api.selectedRoom = room;
  }
}
