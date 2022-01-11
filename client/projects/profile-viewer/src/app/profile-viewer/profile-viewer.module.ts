import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ProfileViewerRoutingModule } from './profile-viewer-routing.module';
import { ProfileViewerComponent, ProfilePictureDialog } from './profile-viewer.component';
import {MatCardModule} from "@angular/material/card";
import {MatDialogModule} from "@angular/material/dialog";
import {MatListModule} from "@angular/material/list";
import {MatButtonModule} from "@angular/material/button";
import {MatIconModule} from "@angular/material/icon";

const EXPORTS = [
  ProfileViewerComponent,
];

@NgModule({
  declarations: [
    ...EXPORTS,
    ProfilePictureDialog
  ],
  imports: [
    CommonModule,
    ProfileViewerRoutingModule,
    MatCardModule,
    MatDialogModule,
    MatListModule,
    MatButtonModule,
    MatIconModule,
  ],
  exports: [...EXPORTS],
})
export class ProfileViewerModule {
  static exports = EXPORTS;
}
