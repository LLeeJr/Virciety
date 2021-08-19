import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {HomeComponent} from "./home/home.component";
import {loadRemoteModule} from "@angular-architects/module-federation";

const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
    pathMatch: 'full',
  },
  {
    path: 'post-comment',
    loadChildren: () =>
      import("./post-comment/post-comment.module").then(m => m.PostCommentModule),
  },
  {
    path: 'mfe1',
    loadChildren: () =>
      loadRemoteModule({
        remoteName: 'mfe1',
        exposedModule: './MfefeatureModule'
      }).then(m => m.MfefeatureModule),
  },
  {
    path: 'post-mfe',
    loadChildren: () =>
      loadRemoteModule({
        remoteName: 'post_mfe',
        exposedModule: './PostModule'
      }).then(m => m.PostModule)
  },
  {
    path: 'comment-mfe',
    loadChildren: () =>
      loadRemoteModule({
        remoteName: 'comment_mfe',
        exposedModule: './CommentModule'
      }).then(m => m.CommentModule)
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
