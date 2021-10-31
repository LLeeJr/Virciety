import { NgModule } from '@angular/core';
import {MatListModule} from "@angular/material/list";
import {MatCardModule} from "@angular/material/card";


@NgModule({
  declarations: [],
  imports: [
    MatListModule,
    MatCardModule,
  ],
  exports: [
    MatListModule,
    MatCardModule,
  ],
})
export class MaterialModule { }
