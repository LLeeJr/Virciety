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
  showSettings = false;
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
        dialogRef.afterClosed().subscribe(room => {
          if (room && room.data && room.data.createRoom) {
            this.chatrooms = [...this.chatrooms, room.data.createRoom];
          }
        });
      }
    });
  }

  removeRoom(room: Room) {
    this.api.deleteRoom(room.name, room.id, room.owner).subscribe(value => {
      if (value && value.data) {
        this.refreshRooms(room.id);
      }
    });
  }

  refreshRooms(id: string) {
    let rooms = [];
    for (let chatroom of this.chatrooms) {
      if (chatroom.id !== id) {
        rooms.push(chatroom);
      }
    }
    this.chatrooms = [...rooms];
    if (sessionStorage.getItem("room")) {
      let room = JSON.parse(<string>sessionStorage.getItem("room"));
      if (room.id === id) {
        sessionStorage.removeItem("room");
      }
    }
  }

  isOwner(room: Room) {
    return room.owner == this.username;
  }

  handleLeaveChatroom(room: Room) {
    if (this.isOwner(room)) {
      let members = [];
      for (let member of room.member) {
        if (member !== this.username) {
          members.push(member)
        }
      }

      let dialogRef = this.dialog.open(SelectOwnerDialog, {
        disableClose: true,
        data: members,
      });

      dialogRef.afterClosed().subscribe((pickedUser) => {
        if (pickedUser) {
          this.api.leaveChat(room.id, this.username, pickedUser).subscribe(value => {
            if (value && value.data && value.data.leaveChat) {
              let { id } = value.data.leaveChat;
              this.refreshRooms(id);
            }
          });
        }
      })
    } else {
      this.api.leaveChat(room.id, this.username, undefined).subscribe(value => {
        if (value && value.data && value.data.leaveChat) {
          let { id } = value.data.leaveChat;
          this.refreshRooms(id);
        }
      });
    }
  }

  parseChatName(roomName: string) {
    let name = roomName;
    let front = this.username.concat('-');
    let back = '-'.concat(this.username);
    if (roomName.includes(front)) {
      name = roomName.replace(front, '')
    } else if (roomName.includes(back)) {
      name = roomName.replace(back, '')
    }
    return name;
  }

  showMembers(room: Room) {
    return !room.isDirect ? `${room.member.length} members` : '';
  }
}

@Component({
  selector: 'add-chat-dialog',
  templateUrl: './add-chat-dialog.html',
  styleUrls: ['./add-chat-dialog.scss']
})
export class AddChatDialog implements OnInit {
  username: string = '';
  roomName: string = '';
  friendList: string[] = [];
  pickedUsers: string[] = [];
  friends = new FormControl();
  // friends = new FormControl([], [
  //   Validators.required,
  //   Validators.minLength(1),
  // ]);
  nameInput = new FormControl('', [
    Validators.required,
    Validators.minLength(2),
  ]);
  @Output() newRoom = new EventEmitter<any>();
  checked = false;

  constructor(@Inject(MAT_DIALOG_DATA) public data: any,
              private api: ApiService,
              private dialogRef: MatDialogRef<AddChatDialog>) {
  }

  ngOnInit() {
    let {follows, followers} = this.data;
    this.friendList = this.removeDuplicates(follows.concat(followers));
    console.log(this.friendList)

    this.friends.valueChanges.subscribe(value => {
      console.log(value)
      this.pickedUsers = value;
      if (this.pickedUsers.length == 1) {
        this.nameInput.setValue(`${this.data.username}-${this.pickedUsers[0]}`);
      } else {
        this.checked = false;
      }
    });
  }

  createRoom(name: string, users: string[], checked: boolean) {
    let member = [this.data.username, ...users];
    this.api.createRoom(member, name, this.data.username, checked).subscribe(value => {
      this.dialogRef.close(value);
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
    return this.nameInput.valid && this.pickedUsers.length > 0;
  }
}

@Component({
  selector: 'select-owner-dialog',
  templateUrl: './select-owner-dialog.html',
  styleUrls: ['./select-owner-dialog.scss']
})
export class SelectOwnerDialog {

  members: string[] = [];
  pickedUser: string;

  constructor(@Inject(MAT_DIALOG_DATA) public data: any,
              private dialogRef: MatDialogRef<SelectOwnerDialog>) {
    this.members = data;
  }

  submit() {
    this.dialogRef.close(this.pickedUser);
  }
}
