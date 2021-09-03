import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {FileUploadTestComponent} from "./file-upload-test/file-upload-test.component";

const routes: Routes = [
  {
    path: '',
    redirectTo: '/file-upload',
    pathMatch: 'full',
  },
  {
    path: 'post-mfe',
    loadChildren: () =>
      import("./post/post.module").then(m => m.PostModule),
  },
  {
    path: 'file-upload',
    component: FileUploadTestComponent,
    pathMatch: 'full',
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
