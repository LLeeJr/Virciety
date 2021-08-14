import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    redirectTo: '/comment-mfe',
    pathMatch: 'full',
  },
  {
    path: 'comment-mfe',
    loadChildren: () =>
      import("./comment/comment.module").then(m => m.CommentModule),
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
