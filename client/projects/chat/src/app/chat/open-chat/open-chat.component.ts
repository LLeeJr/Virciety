import {Component, ElementRef, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from "../../api/api.service";
import {Dm} from "../../data/dm";
import {KeycloakService} from "keycloak-angular";
import {Subscription} from "rxjs";
import {ActivatedRoute, Params, Router} from "@angular/router";
import {Room} from "../../data/room";
import {DatePipe} from "@angular/common";

@Component({
  selector: 'app-open-chat',
  templateUrl: './open-chat.component.html',
  styleUrls: ['./open-chat.component.scss']
})
export class OpenChatComponent implements OnInit, OnDestroy {

  messages: Dm[] = [];

  message: string = '';
  username: string;
  private subscription: Subscription;
  private room: Room;

  constructor(public api: ApiService,
              public datePipe: DatePipe,
              private elem: ElementRef,
              private keycloak: KeycloakService,
              private route: ActivatedRoute,
              private router: Router) { }

  ngOnDestroy(): void {
    this.messages = [];
    this.subscription.unsubscribe();
    this.api.unsubscribeToChat();
    this.api.selectedRoom = null;
  }

  ngOnInit(): void {
    if (this.api.selectedRoom) {
      sessionStorage.setItem("room", JSON.stringify(this.api.selectedRoom));
    }
    this.room = this.getCurrentRoom();
    this.keycloak.isLoggedIn().then(isLoggedIn => {
      if (isLoggedIn) {
        this.username = this.keycloak.getUsername();

        this.route.params.subscribe((params: Params) => {
          let roomName = params['name'];
          let users = roomName.split('-');
          if (this.room && this.room.id) {
            this.api.getRoom(roomName, this.room.id).subscribe(value => {
              if (value && value.data && value.data.getRoom) {
                this.api.selectedRoom = value.data.getRoom;
                let id = value.data.getRoom.id;
                this.getMessagesAndSubscribe(id, roomName);
              }
            });
          } else {
            // id is empty, so request comes from a direct chat call
            this.api.getDirectRoom(users[0], users[1]).subscribe(value => {
              if (value && value.data && value.data.getDirectRoom) {
                this.api.selectedRoom = value.data.getDirectRoom;
                let id = value.data.getDirectRoom.id;
                this.getMessagesAndSubscribe(id, roomName);
              }
            }, () => {
              // room was not found, so a new one needs to be created
              this.api.createRoom(users, roomName, this.username, true).subscribe(value => {
                if (value && value.data && value.data.createRoom) {
                  this.api.selectedRoom = value.data.createRoom;
                  let id = value.data.createRoom.id;
                  this.getMessagesAndSubscribe(id, roomName);
                }
              })
            });
          }
        });
      }
    });
  }

  private getCurrentRoom() {
    return this.api.selectedRoom ? this.api.selectedRoom : JSON.parse(<string>sessionStorage.getItem("room"));
  }

  private getMessagesAndSubscribe(id: string, roomName: string) {
    this.subscription = this.api.getMessagesFromRoom(id).subscribe(value => {
      this.messages = value.data.getMessagesFromRoom;
    });

    this.api.subscribeToChat(roomName).subscribe(value => {
      if (!value || !value.data || !value.data.dmAdded) {
        return;
      }

      const {data} = value

      const {chatroomId, createdAt, createdBy, msg, __typename} = data.dmAdded;

      const dm = new Dm(chatroomId, createdAt, createdBy, msg, __typename);
      this.messages = Object.assign([], this.messages)
      this.messages.push(dm);
    });
  }

  async sendMessage() {
    if (this.message.length > 0) {
      await this.api.writeDm(this.message).toPromise();
      this.message = '';
      this.scrollToBottom();
    }
  }

  private scrollToBottom() {
    try {
      this.elem.nativeElement.scrollTop = this.elem.nativeElement.scrollHeight;
    } catch (err) { }
  }

  goBack() {
    this.router.navigate(['/chat']);
  }

  getRoomMembers() {
    return this.api.selectedRoom.member.join(', ');
  }

  isDirect() {
    return this.room.isDirect;
  }

  transformTime(createdAt: any) {
    return this.datePipe.transform(createdAt, 'short');
  }
}
