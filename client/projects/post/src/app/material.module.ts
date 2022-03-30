import { NgModule } from '@angular/core';
import {MatCardModule} from "@angular/material/card";
import {MatButtonModule} from "@angular/material/button";
import {MatProgressSpinnerModule} from "@angular/material/progress-spinner";
import {MatDialogModule} from "@angular/material/dialog";
import {MatIconModule} from "@angular/material/icon";
import {MatInputModule} from "@angular/material/input";
import {MatListModule} from "@angular/material/list";
import {MatLineModule} from "@angular/material/core";
import {MatTooltipModule} from "@angular/material/tooltip";

@NgModule({
  imports: [
    MatCardModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatLineModule,
    MatTooltipModule
  ],
  exports: [
    MatCardModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatLineModule,
    MatTooltipModule
  ],
})
export class MaterialModule { }
