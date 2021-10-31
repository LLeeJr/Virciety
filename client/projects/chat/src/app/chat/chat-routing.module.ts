import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {ChatComponent} from "./chat.component";
import {OpenChatComponent} from "./open-chat/open-chat.component";

const routes: Routes = [
  {
    path: '',
    component: ChatComponent
  },
  {
    path: ':name',
    component: OpenChatComponent
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ChatRoutingModule { }
