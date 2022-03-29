import {Component, ElementRef, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from "../../api/api.service";
import {Dm} from "../../data/dm";
import {KeycloakService} from "keycloak-angular";
import {Subscription} from "rxjs";
import {ActivatedRoute, Params, Router} from "@angular/router";
import {Room} from "../../data/room";
import {DatePipe} from "@angular/common";
import {take} from "rxjs/operators";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";

@Component({
  selector: 'app-open-chat',
  templateUrl: './open-chat.component.html',
  styleUrls: ['./open-chat.component.scss'],
  exportAs: 'OpenChatComponent',
})
export class OpenChatComponent implements OnInit, OnDestroy {

  messages: Dm[] = [];

  message: string = '';
  username: string;
  room: Room;
  isPhonePortrait: boolean = false;
  private subscription: Subscription;

  constructor(public api: ApiService,
              public datePipe: DatePipe,
              private elem: ElementRef,
              private keycloak: KeycloakService,
              private route: ActivatedRoute,
              private router: Router,
              private responsive: BreakpointObserver) { }

  ngOnDestroy(): void {
    this.messages = [];
    this.subscription.unsubscribe();
    this.api.unsubscribeToChat();
    this.api.selectedRoom = null;
  }

  ngOnInit(): void {
    this.responsive.observe(Breakpoints.HandsetPortrait).subscribe((result) => {
      this.isPhonePortrait = result.matches;
    });

    this.keycloak.isLoggedIn().then(isLoggedIn => {
      if (isLoggedIn) {
        this.username = this.keycloak.getUsername();

        this.route.params.subscribe((params: Params) => {
          let roomName = params['name'];
          let users = roomName.split('-');
          this.room = this.getCurrentRoom(roomName);
          if (this.room && this.room.id) {
            this.api.getRoom(roomName, this.room.id)
              .pipe(take(1))
              .subscribe(value => {
              if (value && value.data && value.data.getRoom) {
                this.room = value.data.getRoom;
                this.storeRoom(this.room);
                this.getMessagesAndSubscribe(this.room.id, roomName);
              }
            });
          } else {
            // id is empty, so request comes from a direct chat call
            this.api.getDirectRoom(users[0], users[1])
              .pipe(take(1))
              .subscribe({
                next: (value) => {
                  if (value && value.data && value.data.getDirectRoom) {
                    this.room = value.data.getDirectRoom;
                    this.storeRoom(this.room);
                    this.getMessagesAndSubscribe(this.room.id, roomName);
                  }
                },
                error: () => {
                  // room was not found, so a new one needs to be created
                  this.api.createRoom(users, roomName, this.username, true)
                    .pipe(take(1))
                    .subscribe(value => {
                      if (value && value.data && value.data.createRoom) {
                        this.room = value.data.createRoom;
                        this.storeRoom(this.room);
                        this.getMessagesAndSubscribe(this.room.id, roomName);
                      }
                    })
                }
              });
          }
        });
      }
    });
  }

  private storeRoom(room: Room) {
    sessionStorage.setItem("room", JSON.stringify(room));
  }

  private getCurrentRoom(roomName: string) {
    if (this.api.selectedRoom) {
      return this.api.selectedRoom;
    }
    if (sessionStorage.getItem("room")) {
      let room =  JSON.parse(<string>sessionStorage.getItem("room"));
      if (room.name === roomName) {
        sessionStorage.removeItem("room");
        return room;
      }
    }
  }

  private getMessagesAndSubscribe(id: string, roomName: string) {
    this.subscription = this.api.getMessagesFromRoom(id).subscribe(value => {
      this.messages = value.data.getMessagesFromRoom;
      this.scrollToBottom();
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
      await this.api.writeDm(this.message, this.room.name, this.room.id, this.username).toPromise();
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
    return this.room.member.join(', ');
  }

  isDirect() {
    if (this.room) {
      return this.room.isDirect;
    }
    return false;
  }

  transformTime(createdAt: any) {
    return this.datePipe.transform(createdAt, 'short');
  }
}
