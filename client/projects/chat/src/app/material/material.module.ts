import { NgModule } from '@angular/core';
import {MatListModule} from "@angular/material/list";
import {MatCardModule} from "@angular/material/card";
import {MatIconModule} from "@angular/material/icon";
import {MatInputModule} from "@angular/material/input";
import {MatButtonModule} from "@angular/material/button";
import {MatDialogModule} from "@angular/material/dialog";
import {MatSelectModule} from "@angular/material/select";
import {MatOptionModule} from "@angular/material/core";
import {MatRadioModule} from "@angular/material/radio";


@NgModule({
  declarations: [],
  imports: [
    MatListModule,
    MatCardModule,
    MatIconModule,
    MatInputModule,
    MatButtonModule,
    MatDialogModule,
    MatSelectModule,
    MatOptionModule,
    MatRadioModule,
  ],
  exports: [
    MatListModule,
    MatCardModule,
    MatIconModule,
    MatInputModule,
    MatButtonModule,
    MatDialogModule,
    MatSelectModule,
    MatOptionModule,
    MatRadioModule,
  ],
})
export class MaterialModule { }
