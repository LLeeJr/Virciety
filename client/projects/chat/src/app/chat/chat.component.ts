import {Component, EventEmitter, Inject, OnInit, Output} from '@angular/core';
import {ApiService} from "../api/api.service";
import {AuthLibService} from "auth-lib";
import {Room} from "../data/room";
import {KeycloakService} from "keycloak-angular";
import {Router} from "@angular/router";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {FormControl, Validators} from "@angular/forms";

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {

  chatrooms: Room[] = [];
  private username: string;

  constructor(private api: ApiService,
              private auth: AuthLibService,
              private dialog: MatDialog,
              private keycloak: KeycloakService,
              private router: Router) {
  }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();
          this.api.getRoomsByUser(this.username).subscribe(value => {
            this.chatrooms = value.data.getRoomsByUser;
          });
        })
      } else {
        this.keycloak.login();
      }
    });
  }

  setChatPartner(room: Room) {
    this.api.chatMembers = room.member;
    this.api.selectedRoom = room;
    this.router.navigate([`chat/${room.name}`]);
  }

  openDialog() {
    this.auth.getUserByName(this.username).subscribe(value => {
      if (value && value.data && value.data.getUserByName) {
        let dialogRef = this.dialog.open(AddChatDialog, {
          data: value.data.getUserByName,
        });
        dialogRef.componentInstance.newRoom.subscribe(room => {
          this.chatrooms = [...this.chatrooms, room];
        });
      }
    });
  }

  removeRoom(room: Room) {
    this.api.deleteRoom(room.name, room.id, room.owner).subscribe(value => {
      if (value && value.data) {
        let rooms = [];
        for (let chatroom of this.chatrooms) {
          if (chatroom.id !== room.id) {
            rooms.push(chatroom);
          }
        }
        this.chatrooms = [...rooms];
      }
    });
  }

  isOwner(room: Room) {
    return room.owner == this.username;
  }
}

@Component({
  selector: 'add-chat-dialog',
  templateUrl: './add-chat-dialog.html',
  styleUrls: ['./add-chat-dialog.scss']
})
export class AddChatDialog {
  username: string = '';
  roomName: string = '';
  friendList: string[] = [];
  pickedUsers: string[] = [];
  friends = new FormControl([], [
    Validators.required,
    Validators.minLength(2),
  ]);
  nameInput = new FormControl('', [
    Validators.required,
    Validators.minLength(2),
  ]);
  @Output() newRoom = new EventEmitter<any>();

  constructor(@Inject(MAT_DIALOG_DATA) public data: any,
              private api: ApiService,
              private dialogRef: MatDialogRef<AddChatDialog>) {
    if (data) {
      let {follows, followers} = this.data;
      this.friendList = this.removeDuplicates(follows.concat(followers));
    }

    this.friends.valueChanges.subscribe(value => this.pickedUsers = value);
  }

  createRoom(name: string, users: string[]) {
    let member = [this.data.username, ...users];
    this.api.createRoom(member, name, this.data.username).subscribe(value => {
      if (value && value.data && value.data.createRoom) {
        this.newRoom.emit(value.data.createRoom);
        this.dialogRef.close();
      }
    });
  }

  removeDuplicates(list: string[]) {
    let l = list.concat();
    for (let i = 0; i < l.length; i++) {
      for (let j = i+1; j < l.length; j++) {
        if (l[i] === l[j]) {
          l.splice(j--, 1);
        }
      }
    }
    return l;
  }

  valid() {
    return this.nameInput.valid && this.pickedUsers.length > 1;
  }
}
