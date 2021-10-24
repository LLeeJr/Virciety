import { Component, OnInit } from '@angular/core';
import {ApiService} from "../api/api.service";
import {AuthLibService} from "auth-lib";
import {Room} from "../data/room";
import {KeycloakService} from "keycloak-angular";
import {Router} from "@angular/router";

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {

  chatrooms: Room[] = [];

  constructor(private api: ApiService,
              private auth: AuthLibService,
              private keycloak: KeycloakService,
              private router: Router) {
  }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        let username = this.keycloak.getUsername();
        this.api.getRoomsByUser(username).subscribe(value => {
          this.chatrooms = value.data.getRoomsByUser;
        });
      }
    });
  }

  setChatPartner(room: Room) {
    this.api.chatMembers = room.member;
    this.api.selectedRoom = room;
    this.router.navigate([`chat/${room.name}`]);
  }
}
