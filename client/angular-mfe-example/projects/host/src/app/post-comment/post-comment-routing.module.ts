import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {PostCommentComponent} from "./post-comment.component";

const routes: Routes = [
  {
    path: '',
    component: PostCommentComponent
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class PostCommentRoutingModule { }
