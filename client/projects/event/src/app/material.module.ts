import { NgModule } from '@angular/core';
import {MatDatepickerModule} from "@angular/material/datepicker";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatInputModule} from "@angular/material/input";
import {MatNativeDateModule} from "@angular/material/core";
import {MatCardModule} from "@angular/material/card";
import {MatDividerModule} from "@angular/material/divider";
import {MatIconModule} from "@angular/material/icon";
import {MatCheckboxModule} from "@angular/material/checkbox";
import {NgxMatTimepickerModule} from "ngx-mat-timepicker";
import {MatDialogModule} from "@angular/material/dialog";
import {MatButtonModule} from "@angular/material/button";
import {MatToolbarModule} from "@angular/material/toolbar";
import {MatListModule} from "@angular/material/list";

@NgModule({
  imports: [
    MatCardModule,
    MatDatepickerModule,
    MatFormFieldModule,
    MatInputModule,
    MatNativeDateModule,
    MatDividerModule,
    MatIconModule,
    MatCheckboxModule,
    NgxMatTimepickerModule,
    MatDialogModule,
    MatButtonModule,
    MatToolbarModule,
    MatListModule
  ],
  exports: [
    MatCardModule,
    MatDatepickerModule,
    MatFormFieldModule,
    MatInputModule,
    MatNativeDateModule,
    MatDividerModule,
    MatIconModule,
    MatCheckboxModule,
    NgxMatTimepickerModule,
    MatDialogModule,
    MatButtonModule,
    MatToolbarModule,
    MatListModule
  ],
  providers: [
    MatDatepickerModule
  ]
})
export class MaterialModule { }
