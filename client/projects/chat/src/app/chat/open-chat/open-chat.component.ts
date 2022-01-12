import {AfterViewInit, Component, ElementRef, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {ApiService} from "../../api/api.service";
import {Dm} from "../../data/dm";
import {KeycloakService} from "keycloak-angular";
import {Subscription} from "rxjs";
import {ActivatedRoute, Params} from "@angular/router";

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

  constructor(public api: ApiService,
              private elem: ElementRef,
              private keycloak: KeycloakService,
              private route: ActivatedRoute) { }

  ngOnDestroy(): void {
    this.messages = [];
    this.subscription.unsubscribe();
    this.api.unsubscribeToChat();
    this.api.selectedRoom = null;
  }

  ngOnInit(): void {
    this.keycloak.isLoggedIn().then(isLoggedIn => {
      if (isLoggedIn) {
        this.username = this.keycloak.getUsername();
      }
    });

    this.route.params.subscribe((params: Params) => {
      let roomName = params['name'];
      this.subscription = this.api.getMessagesFromRoom(roomName).subscribe(value => {
        this.messages = value.data.getMessagesFromRoom;
      });

      this.api.subscribeToChat(roomName).subscribe(value => {
        if (!value || !value.data || !value.data.dmAdded) {
          return;
        }

        const { data } = value

        const { chatroomId, createdAt, createdBy, msg, __typename } = data.dmAdded;

        const dm = new Dm(chatroomId, createdAt, createdBy, msg, __typename);
        this.messages = Object.assign([], this.messages)
        this.messages.push(dm);
      });
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
}
