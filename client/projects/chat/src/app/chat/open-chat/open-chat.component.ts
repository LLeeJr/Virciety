import {Component, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from "../../api/api.service";
import {Dm} from "../../data/dm";

@Component({
  selector: 'app-open-chat',
  templateUrl: './open-chat.component.html',
  styleUrls: ['./open-chat.component.scss']
})
export class OpenChatComponent implements OnInit, OnDestroy {

  user1: string | null = '';
  messages: Dm[] = [];

  message: string = '';

  constructor(public api: ApiService) { }

  ngOnDestroy(): void {
    this.messages = [];
  }

  ngOnInit(): void {
    if (this.user1 !== null && this.api.chatMembers.length !== 0) {
      this.api.getMessagesFromRoom().subscribe(value => {
        this.messages = value.data.getMessagesFromRoom;
      });
    }

    this.api.subscribeToChat().subscribe(value => {
      if (!value || !value.data || !value.data.dmAdded) {
        return;
      }

      const { data } = value

      const { chatroomId, createdAt, createdBy, msg, __typename } = data.dmAdded;

      const dm = new Dm(chatroomId, createdAt, createdBy, msg, __typename);
      this.messages = Object.assign([], this.messages)
      this.messages.push(dm);
    });

  }

  async sendMessage() {
    if (this.message.length > 0) {
      await this.api.writeDm(this.message).toPromise();
      this.message = '';
    }
  }

}
