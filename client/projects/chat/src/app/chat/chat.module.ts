import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ChatRoutingModule } from './chat-routing.module';
import { ChatComponent } from './chat.component';
import {Apollo} from "apollo-angular";
import {InMemoryCache} from "@apollo/client/core";
import {HttpLink} from "apollo-angular/http";
import { OpenChatComponent } from './open-chat/open-chat.component';
import {FormsModule} from "@angular/forms";


@NgModule({
  declarations: [
    ChatComponent,
    OpenChatComponent
  ],
  imports: [
    CommonModule,
    ChatRoutingModule,
    FormsModule
  ],
})
export class ChatModule {
  constructor(
    apollo: Apollo,
    httpLink: HttpLink
  ) {
    const link = httpLink.create({
      uri: 'http://localhost:8081/query'
    });

    apollo.create({
      link: link,
      cache: new InMemoryCache(),
    });
  }
}
