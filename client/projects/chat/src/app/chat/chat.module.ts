import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ChatRoutingModule } from './chat-routing.module';
import { ChatComponent } from './chat.component';
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
export class ChatModule { }