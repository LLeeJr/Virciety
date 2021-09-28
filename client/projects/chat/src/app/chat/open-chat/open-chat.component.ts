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
    console.log('Open Chat');

    if (this.user1 !== null && this.api.chatPartner !== null) {
      this.api.getChat().subscribe(value => {
        this.messages = value.data.getChat;
      });
    }

    console.log('Subscribe to chat...');
    this.api.subscribeToChat().subscribe(value => {
      if (!value || !value.data || !value.data.dmAdded) {
        return;
      }

      const { data } = value

      const {__typename, id, msg} = data.dmAdded;

      const dm = new Dm(__typename, id, msg);
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

  getDate(message: { __typename: string; id: string; msg: string }) {
    const data = message.id.split('__');
    return data[1];
  }
}
