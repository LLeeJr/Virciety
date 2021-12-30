import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ProfileViewerRoutingModule } from './profile-viewer-routing.module';
import { ProfileViewerComponent } from './profile-viewer.component';


@NgModule({
  declarations: [
    ProfileViewerComponent
  ],
  imports: [
    CommonModule,
    ProfileViewerRoutingModule
  ]
})
export class ProfileViewerModule { }
