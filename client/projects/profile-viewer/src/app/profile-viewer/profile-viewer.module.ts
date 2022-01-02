import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ProfileViewerRoutingModule } from './profile-viewer-routing.module';
import { ProfileViewerComponent, ProfilePictureDialog } from './profile-viewer.component';
import {MatCardModule} from "@angular/material/card";
import {MatDialogModule} from "@angular/material/dialog";
import {MatListModule} from "@angular/material/list";
import {MatButtonModule} from "@angular/material/button";


@NgModule({
  declarations: [
    ProfileViewerComponent,
    ProfilePictureDialog
  ],
  imports: [
    CommonModule,
    ProfileViewerRoutingModule,
    MatCardModule,
    MatDialogModule,
    MatListModule,
    MatButtonModule,
  ]
})
export class ProfileViewerModule { }
