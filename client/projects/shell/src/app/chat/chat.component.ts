import { Component, OnInit } from '@angular/core';
import {Router} from "@angular/router";

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {

  constructor(private router: Router) { }

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
