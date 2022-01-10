import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {DialogCreatePostComponent} from "./create-post.component";

const routes: Routes = [
  {
    path: '',
    component: DialogCreatePostComponent,
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CreatePostRoutingModule { }
