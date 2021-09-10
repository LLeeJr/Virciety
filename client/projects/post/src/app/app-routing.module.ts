import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    redirectTo: '/create_post',
    pathMatch: 'full',
  },
  {
    path: 'post',
    loadChildren: () =>
      import("./post/post.module").then(m => m.PostModule),
  },
  {
    path: 'create_post',
    loadChildren: () =>
      import("./create-post/create-post.module").then(m => m.CreatePostModule),
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
