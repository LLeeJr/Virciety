import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ProfileViewerRoutingModule } from './profile-viewer-routing.module';
import { ProfileViewerComponent } from './profile-viewer.component';
import {MatCardModule} from "@angular/material/card";


@NgModule({
  declarations: [
    ProfileViewerComponent
  ],
    imports: [
        CommonModule,
        ProfileViewerRoutingModule,
        MatCardModule
    ]
})
export class ProfileViewerModule { }
