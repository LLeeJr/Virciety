import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ChatRoutingModule } from './chat-routing.module';
import {AddChatDialog, ChatComponent} from './chat.component';
import { OpenChatComponent } from './open-chat/open-chat.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {MaterialModule} from "../material/material.module";


@NgModule({
  declarations: [
    ChatComponent,
    OpenChatComponent,
    AddChatDialog,
  ],
  imports: [
    CommonModule,
    ChatRoutingModule,
    FormsModule,
    MaterialModule,
    ReactiveFormsModule,
  ],
})
export class ChatModule { }
