import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {ApiService} from "../../api/api.service";
import {AuthLibService} from "auth-lib";

@Component({
  selector: 'app-open-chat',
  templateUrl: './open-chat.component.html',
  styleUrls: ['./open-chat.component.scss']
})
export class OpenChatComponent implements OnInit {

  user1: string | null = '';
  user2: string | null = '';
  messages: {
    __typename: string,
    id: string,
    msg: string,
  }[] = [];

  message: string = '';

  constructor(private readonly route: ActivatedRoute,
              private api: ApiService,
              private auth: AuthLibService) { }

  ngOnInit(): void {

    this.route.paramMap.subscribe((params) => {
      if (params.has('id')) {
        this.user2 = this.route.snapshot.paramMap.get('id');
      }
    })

    this.user1 = this.auth.userName;

    if (this.user1 !== null && this.user2 !== null) {
      this.api.getChat(this.user1, this.user2).subscribe(value => {
        this.messages = value.data.getChat;
      });
    }

  }

  async sendMessage() {
    if (this.user1 !== null && this.user2 !== null && this.message.length > 0) {
      await this.api.writeDm(this.message, this.user1, this.user2).toPromise();
      this.message = '';
    }
  }

  getDate(message: { __typename: string; id: string; msg: string }) {
    const data = message.id.split('__');
    return data[1];
  }
}
