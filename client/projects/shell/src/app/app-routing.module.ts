import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {AuthGuard} from "./guard/auth.guard";
import {ProfileComponent} from "./profile/profile.component";
import {SinglePostComponent} from "./single-post/single-post.component";
import {SingleEventComponent} from "./single-event/single-event.component";
import {PostComponent} from "./post/post.component";
import {EventComponent} from "./event/event.component";
import {ChatComponent} from "./chat/chat.component";
import {OpenChatComponent} from "./open-chat/open-chat.component";

const routes: Routes = [
  {
    path: '',
    redirectTo: '/home',
    pathMatch: 'full'
  },
  {
    path: 'chat',
    canActivate: [AuthGuard],
    component: ChatComponent,
  },
  {
    path: 'chat/:name',
    canActivate: [AuthGuard],
    component: OpenChatComponent,
  },
  {
    path: 'home',
    canActivate: [AuthGuard],
    component: PostComponent
  },
  {
    path: 'event',
    canActivate: [AuthGuard],
    component: EventComponent,
  },
  {
    path: 'profile',
    canActivate: [AuthGuard],
    component: ProfileComponent
  },
  {
    path: 'p/:id',
    canActivate: [AuthGuard],
    component: SinglePostComponent
  },
  {
    path: 'e/:id',
    canActivate: [AuthGuard],
    component: SingleEventComponent
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
