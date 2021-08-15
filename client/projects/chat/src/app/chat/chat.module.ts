import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ChatRoutingModule } from './chat-routing.module';
import { ChatComponent } from './chat.component';
import {Apollo} from "apollo-angular";
import {InMemoryCache} from "@apollo/client/core";
import {HttpLink} from "apollo-angular/http";


@NgModule({
  declarations: [
    ChatComponent
  ],
  imports: [
    CommonModule,
    ChatRoutingModule
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
