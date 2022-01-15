import { NgModule } from '@angular/core';
import {MatDatepickerModule} from "@angular/material/datepicker";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatInputModule} from "@angular/material/input";
import {MatNativeDateModule} from "@angular/material/core";
import {MatCardModule} from "@angular/material/card";
import {MatDividerModule} from "@angular/material/divider";

@NgModule({
  imports: [
    MatCardModule,
    MatDatepickerModule,
    MatFormFieldModule,
    MatInputModule,
    MatNativeDateModule,
    MatDividerModule
  ],
  exports: [
    MatCardModule,
    MatDatepickerModule,
    MatFormFieldModule,
    MatInputModule,
    MatNativeDateModule,
    MatDividerModule
  ],
  providers: [
    MatDatepickerModule
  ]
})
export class MaterialModule { }
