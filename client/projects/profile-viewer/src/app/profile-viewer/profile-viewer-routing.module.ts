import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {ProfileViewerComponent} from "./profile-viewer.component";

const routes: Routes = [
  {
    path: '',
    component: ProfileViewerComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ProfileViewerRoutingModule { }
