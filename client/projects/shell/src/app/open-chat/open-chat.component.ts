import { Component, OnInit } from '@angular/core';
import {Router} from "@angular/router";
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-open-chat',
  templateUrl: './open-chat.component.html',
  styleUrls: ['./open-chat.component.scss']
})
export class OpenChatComponent implements OnInit {
  chatMFE: string;

  constructor(private router: Router) {
    this.chatMFE = environment.chatMFE;
  }

  ngOnInit(): void {
  }

  handleError(event: any) {
    let {error, component} = event;
    if (error && component === 'chat') {
      let msg = `${component} is currently offline!`;
      this.router.navigate(['/page-not-found', msg])
    }
  }
}
