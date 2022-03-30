import { Component, OnInit } from '@angular/core';
import {environment} from "../../environments/environment";
import {Router} from "@angular/router";

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {
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
