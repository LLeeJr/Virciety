import { NgModule } from '@angular/core';
import {MatListModule} from "@angular/material/list";
import {MatCardModule} from "@angular/material/card";
import {MatIconModule} from "@angular/material/icon";
import {MatInputModule} from "@angular/material/input";
import {MatButtonModule} from "@angular/material/button";


@NgModule({
  declarations: [],
  imports: [
    MatListModule,
    MatCardModule,
    MatIconModule,
    MatInputModule,
    MatButtonModule,
  ],
  exports: [
    MatListModule,
    MatCardModule,
    MatIconModule,
    MatInputModule,
    MatButtonModule,
  ],
})
export class MaterialModule { }
