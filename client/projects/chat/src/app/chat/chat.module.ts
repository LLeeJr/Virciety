import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ChatRoutingModule } from './chat-routing.module';
import {AddChatDialog, ChatComponent, SelectOwnerDialog} from './chat.component';
import { OpenChatComponent } from './open-chat/open-chat.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {MaterialModule} from "../material/material.module";
import {MatFormFieldModule} from "@angular/material/form-field";


const EXPORTS = [
  ChatComponent,
  OpenChatComponent
]

@NgModule({
  declarations: [
    ...EXPORTS,
    AddChatDialog,
    SelectOwnerDialog,
  ],
  imports: [
    CommonModule,
    ChatRoutingModule,
    FormsModule,
    MaterialModule,
    ReactiveFormsModule,
    MatFormFieldModule,
  ],
  exports: [...EXPORTS]
})
export class ChatModule {
  static exports = EXPORTS
}
